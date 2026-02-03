package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

var (
	dataPath    = flag.String("datapath", "./data", "Path to data files")
	outputDir   = flag.String("outputdir", "./", "Output directory")
	domainLists = flag.String("domainlists", "", "Comma-separated domain list names for geosite.dat")
	ipLists     = flag.String("iplists", "", "Comma-separated IP list names for geoip.dat")
	exportLists = flag.String("exportlists", "", "Comma-separated list names to export as plaintext")
)

// DomainEntry represents a parsed domain entry with type and attributes
type DomainEntry struct {
	Type       routercommon.Domain_Type
	Value      string
	Attributes []string
}

// parseEntry parses a single line into a DomainEntry
// Supports: full:, regexp:, keyword:, domain:, and plain domains
// Also supports attributes: domain.tld@attr1@attr2
func parseEntry(line string) (*DomainEntry, string, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, "", nil
	}

	// Check for include directive
	if strings.HasPrefix(line, "include:") {
		return nil, strings.TrimPrefix(line, "include:"), nil
	}

	entry := &DomainEntry{
		Type: routercommon.Domain_RootDomain,
	}

	// Parse attributes (@attr1@attr2)
	if atIdx := strings.Index(line, "@"); atIdx != -1 {
		attrPart := line[atIdx+1:]
		line = line[:atIdx]
		entry.Attributes = strings.Split(attrPart, "@")
	}

	// Remove trailing whitespace and inline comments
	if spaceIdx := strings.Index(line, " "); spaceIdx != -1 {
		line = line[:spaceIdx]
	}

	// Parse type prefix
	switch {
	case strings.HasPrefix(line, "full:"):
		entry.Type = routercommon.Domain_Full
		entry.Value = strings.TrimPrefix(line, "full:")
	case strings.HasPrefix(line, "regexp:"):
		entry.Type = routercommon.Domain_Regex
		entry.Value = strings.TrimPrefix(line, "regexp:")
	case strings.HasPrefix(line, "keyword:"):
		entry.Type = routercommon.Domain_Plain
		entry.Value = strings.TrimPrefix(line, "keyword:")
	case strings.HasPrefix(line, "domain:"):
		entry.Type = routercommon.Domain_RootDomain
		entry.Value = strings.TrimPrefix(line, "domain:")
	default:
		entry.Type = routercommon.Domain_RootDomain
		entry.Value = line
	}

	if entry.Value == "" {
		return nil, "", nil
	}

	return entry, "", nil
}

// loadDomainList loads a domain list file and recursively resolves includes
func loadDomainList(name string, basePath string, loaded map[string]bool) ([]*DomainEntry, error) {
	// Prevent infinite recursion
	if loaded[name] {
		return nil, nil
	}
	loaded[name] = true

	path := filepath.Join(basePath, name)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	var entries []*DomainEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		entry, includeName, err := parseEntry(line)
		if err != nil {
			fmt.Printf("Warning: %s: %v\n", name, err)
			continue
		}

		if includeName != "" {
			// Handle include directive
			// Parse include attributes: include:list@attr
			includeAttr := ""
			if atIdx := strings.Index(includeName, "@"); atIdx != -1 {
				includeAttr = includeName[atIdx+1:]
				includeName = includeName[:atIdx]
			}

			included, err := loadDomainList(includeName, basePath, loaded)
			if err != nil {
				fmt.Printf("Warning: cannot include %s: %v\n", includeName, err)
				continue
			}

			// Filter by attribute if specified
			for _, inc := range included {
				if includeAttr != "" {
					// Only include entries that have the specified attribute
					hasAttr := false
					for _, attr := range inc.Attributes {
						if attr == includeAttr {
							hasAttr = true
							break
						}
					}
					if !hasAttr {
						continue
					}
				}
				entries = append(entries, inc)
			}
		} else if entry != nil {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// loadIPList loads an IP list file (CIDR format, one per line)
func loadIPList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entries = append(entries, line)
	}
	return entries, scanner.Err()
}

// buildGeoSite creates a GeoSite from domain entries
func buildGeoSite(name string, entries []*DomainEntry) *routercommon.GeoSite {
	site := &routercommon.GeoSite{
		CountryCode: strings.ToUpper(strings.ReplaceAll(name, "-", "_")),
	}

	// Deduplicate entries
	seen := make(map[string]bool)
	for _, e := range entries {
		key := fmt.Sprintf("%d:%s", e.Type, e.Value)
		if seen[key] {
			continue
		}
		seen[key] = true

		domain := &routercommon.Domain{
			Type:  e.Type,
			Value: e.Value,
		}

		// Add attributes
		for _, attr := range e.Attributes {
			if attr == "" {
				continue
			}
			parts := strings.SplitN(attr, "=", 2)
			a := &routercommon.Domain_Attribute{Key: parts[0]}
			if len(parts) == 2 {
				a.TypedValue = &routercommon.Domain_Attribute_BoolValue{BoolValue: parts[1] == "true"}
			} else {
				a.TypedValue = &routercommon.Domain_Attribute_BoolValue{BoolValue: true}
			}
			domain.Attribute = append(domain.Attribute, a)
		}

		site.Domain = append(site.Domain, domain)
	}

	return site
}

// buildGeoIP creates a GeoIP from CIDR entries (supports both IPv4 and IPv6)
func buildGeoIP(name string, cidrs []string) *routercommon.GeoIP {
	geoip := &routercommon.GeoIP{
		CountryCode: strings.ToUpper(strings.ReplaceAll(name, "-", "_")),
	}

	seen := make(map[string]bool)
	for _, cidrStr := range cidrs {
		cidrStr = strings.TrimSpace(cidrStr)

		// Add default mask if not specified
		if !strings.Contains(cidrStr, "/") {
			if strings.Contains(cidrStr, ":") {
				cidrStr = cidrStr + "/128" // IPv6
			} else {
				cidrStr = cidrStr + "/32" // IPv4
			}
		}

		// Deduplicate
		if seen[cidrStr] {
			continue
		}
		seen[cidrStr] = true

		_, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			fmt.Printf("Warning: invalid CIDR %s: %v\n", cidrStr, err)
			continue
		}

		ones, _ := ipNet.Mask.Size()

		// Handle both IPv4 and IPv6
		ip := ipNet.IP
		if v4 := ip.To4(); v4 != nil {
			ip = v4
		}

		geoip.Cidr = append(geoip.Cidr, &routercommon.CIDR{
			Ip:     ip,
			Prefix: uint32(ones),
		})
	}

	return geoip
}

// exportDomainList exports entries to plaintext format
func exportDomainList(entries []*DomainEntry, path string) error {
	var lines []string
	for _, e := range entries {
		var line string
		switch e.Type {
		case routercommon.Domain_Full:
			line = "full:" + e.Value
		case routercommon.Domain_Regex:
			line = "regexp:" + e.Value
		case routercommon.Domain_Plain:
			line = "keyword:" + e.Value
		default:
			line = e.Value
		}
		if len(e.Attributes) > 0 {
			line += "@" + strings.Join(e.Attributes, "@")
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func main() {
	flag.Parse()

	// Build geosite.dat from domain lists
	if *domainLists != "" {
		geoSiteList := &routercommon.GeoSiteList{}

		for _, name := range strings.Split(*domainLists, ",") {
			name = strings.TrimSpace(name)
			loaded := make(map[string]bool)

			entries, err := loadDomainList(name, *dataPath, loaded)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", name, err)
				continue
			}

			geoSite := buildGeoSite(name, entries)
			geoSiteList.Entry = append(geoSiteList.Entry, geoSite)
			fmt.Printf("GeoSite: %s (%d domains)\n", name, len(geoSite.Domain))
		}

		// Sort by country code
		sort.Slice(geoSiteList.Entry, func(i, j int) bool {
			return geoSiteList.Entry[i].CountryCode < geoSiteList.Entry[j].CountryCode
		})

		data, err := proto.Marshal(geoSiteList)
		if err != nil {
			fmt.Printf("Error marshaling geosite: %v\n", err)
			os.Exit(1)
		}

		outputPath := filepath.Join(*outputDir, "geosite.dat")
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Printf("Error writing geosite.dat: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s (%d bytes)\n", outputPath, len(data))
	}

	// Build geoip.dat from IP lists
	if *ipLists != "" {
		geoIPList := &routercommon.GeoIPList{}

		for _, name := range strings.Split(*ipLists, ",") {
			name = strings.TrimSpace(name)
			path := filepath.Join(*dataPath, name)

			cidrs, err := loadIPList(path)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", name, err)
				continue
			}

			geoIP := buildGeoIP(name, cidrs)
			geoIPList.Entry = append(geoIPList.Entry, geoIP)
			fmt.Printf("GeoIP: %s (%d CIDRs)\n", name, len(geoIP.Cidr))
		}

		// Sort by country code
		sort.Slice(geoIPList.Entry, func(i, j int) bool {
			return geoIPList.Entry[i].CountryCode < geoIPList.Entry[j].CountryCode
		})

		data, err := proto.Marshal(geoIPList)
		if err != nil {
			fmt.Printf("Error marshaling geoip: %v\n", err)
			os.Exit(1)
		}

		outputPath := filepath.Join(*outputDir, "geoip.dat")
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Printf("Error writing geoip.dat: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s (%d bytes)\n", outputPath, len(data))
	}

	// Export plaintext
	if *exportLists != "" {
		for _, name := range strings.Split(*exportLists, ",") {
			name = strings.TrimSpace(name)
			loaded := make(map[string]bool)

			entries, err := loadDomainList(name, *dataPath, loaded)
			if err != nil {
				// Try as IP list
				path := filepath.Join(*dataPath, name)
				ipEntries, ipErr := loadIPList(path)
				if ipErr != nil {
					fmt.Printf("Warning: cannot export %s: %v\n", name, err)
					continue
				}
				// Export IP list as plaintext
				txtPath := filepath.Join(*outputDir, name+".txt")
				sort.Strings(ipEntries)
				if err := os.WriteFile(txtPath, []byte(strings.Join(ipEntries, "\n")+"\n"), 0644); err != nil {
					fmt.Printf("Warning: cannot write %s: %v\n", txtPath, err)
					continue
				}
				fmt.Printf("Exported: %s (%d entries)\n", txtPath, len(ipEntries))
				continue
			}

			txtPath := filepath.Join(*outputDir, name+".txt")
			if err := exportDomainList(entries, txtPath); err != nil {
				fmt.Printf("Warning: cannot write %s: %v\n", txtPath, err)
				continue
			}
			fmt.Printf("Exported: %s (%d entries)\n", txtPath, len(entries))
		}
	}
}
