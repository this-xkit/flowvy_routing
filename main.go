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
	dataPath      = flag.String("datapath", "./data", "Path to data files")
	outputDir     = flag.String("outputdir", "./", "Output directory")
	domainLists   = flag.String("domainlists", "", "Comma-separated domain list names for geosite.dat")
	ipLists       = flag.String("iplists", "", "Comma-separated IP list names for geoip.dat")
	exportLists   = flag.String("exportlists", "", "Comma-separated list names to export as plaintext")
)

func loadList(path string) ([]string, error) {
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

func buildGeoSite(name string, domains []string) *routercommon.GeoSite {
	site := &routercommon.GeoSite{
		CountryCode: strings.ToUpper(strings.ReplaceAll(name, "-", "_")),
	}

	for _, d := range domains {
		domain := &routercommon.Domain{
			Value: d,
			Type:  routercommon.Domain_RootDomain,
		}
		site.Domain = append(site.Domain, domain)
	}

	return site
}

func buildGeoIP(name string, cidrs []string) *routercommon.GeoIP {
	geoip := &routercommon.GeoIP{
		CountryCode: strings.ToUpper(strings.ReplaceAll(name, "-", "_")),
	}

	for _, cidrStr := range cidrs {
		// Add /32 if no mask specified
		if !strings.Contains(cidrStr, "/") {
			cidrStr = cidrStr + "/32"
		}

		_, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			fmt.Printf("Warning: invalid CIDR %s: %v\n", cidrStr, err)
			continue
		}

		ones, _ := ipNet.Mask.Size()
		geoip.Cidr = append(geoip.Cidr, &routercommon.CIDR{
			Ip:     ipNet.IP.To4(),
			Prefix: uint32(ones),
		})
	}

	return geoip
}

func main() {
	flag.Parse()

	// Build geosite.dat from domain lists
	if *domainLists != "" {
		geoSiteList := &routercommon.GeoSiteList{}

		for _, name := range strings.Split(*domainLists, ",") {
			name = strings.TrimSpace(name)
			path := filepath.Join(*dataPath, name)

			domains, err := loadList(path)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", name, err)
				continue
			}

			geoSite := buildGeoSite(name, domains)
			geoSiteList.Entry = append(geoSiteList.Entry, geoSite)
			fmt.Printf("GeoSite: %s (%d domains)\n", name, len(domains))
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

			cidrs, err := loadList(path)
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
			path := filepath.Join(*dataPath, name)

			entries, err := loadList(path)
			if err != nil {
				fmt.Printf("Warning: cannot export %s: %v\n", name, err)
				continue
			}

			sort.Strings(entries)
			txtPath := filepath.Join(*outputDir, name+".txt")
			if err := os.WriteFile(txtPath, []byte(strings.Join(entries, "\n")), 0644); err != nil {
				fmt.Printf("Warning: cannot write %s: %v\n", txtPath, err)
				continue
			}
			fmt.Printf("Exported: %s\n", txtPath)
		}
	}
}
