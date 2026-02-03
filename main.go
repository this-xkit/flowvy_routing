package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

var (
	dataPath    = flag.String("datapath", "./data", "Path to domain list files")
	outputName  = flag.String("outputname", "geosite.dat", "Output filename")
	outputDir   = flag.String("outputdir", "./", "Output directory")
	exportLists = flag.String("exportlists", "", "Comma-separated list names to export as plaintext")
)

type Entry struct {
	Type  string
	Value string
	Attrs []string
}

type List struct {
	Name    string
	Entries []Entry
}

type ParsedList struct {
	Name      string
	Inclusion map[string]bool
	Entries   []Entry
}

var loadedLists = make(map[string]*List)

func Load(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		list := &List{Name: name}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			entry := parseEntry(line)
			list.Entries = append(list.Entries, entry)
		}

		loadedLists[name] = list
		fmt.Printf("Loaded: %s (%d entries)\n", name, len(list.Entries))
		return nil
	})
}

func parseEntry(line string) Entry {
	entry := Entry{Type: "domain"}

	// Split by @ for attributes
	parts := strings.Split(line, "@")
	main := strings.TrimSpace(parts[0])
	for i := 1; i < len(parts); i++ {
		attr := strings.TrimSpace(parts[i])
		if attr != "" {
			entry.Attrs = append(entry.Attrs, attr)
		}
	}

	// Check for type prefix
	if idx := strings.Index(main, ":"); idx != -1 {
		prefix := main[:idx]
		switch prefix {
		case "domain", "full", "keyword", "regexp", "include":
			entry.Type = prefix
			main = main[idx+1:]
		}
	}

	entry.Value = main
	return entry
}

func ParseList(name string, attrFilter map[string]bool) (*ParsedList, error) {
	list, ok := loadedLists[name]
	if !ok {
		return nil, fmt.Errorf("list not found: %s", name)
	}

	parsed := &ParsedList{
		Name:      name,
		Inclusion: make(map[string]bool),
	}

	for _, entry := range list.Entries {
		// Check attribute filter
		if len(attrFilter) > 0 {
			match := false
			for _, attr := range entry.Attrs {
				if attrFilter[attr] {
					match = true
					break
				}
			}
			if !match && len(entry.Attrs) > 0 {
				continue
			}
		}

		if entry.Type == "include" {
			// Parse include directive
			incName := entry.Value
			incAttrs := make(map[string]bool)
			for _, attr := range entry.Attrs {
				incAttrs[attr] = true
			}

			incList, err := ParseList(incName, incAttrs)
			if err != nil {
				fmt.Printf("Warning: failed to include %s: %v\n", incName, err)
				continue
			}

			parsed.Inclusion[incName] = true
			for k := range incList.Inclusion {
				parsed.Inclusion[k] = true
			}
			parsed.Entries = append(parsed.Entries, incList.Entries...)
		} else {
			parsed.Entries = append(parsed.Entries, entry)
		}
	}

	return parsed, nil
}

func toProto(list *ParsedList) *routercommon.GeoSite {
	site := &routercommon.GeoSite{
		CountryCode: strings.ToUpper(list.Name),
	}

	for _, entry := range list.Entries {
		domain := &routercommon.Domain{Value: entry.Value}

		switch entry.Type {
		case "domain":
			domain.Type = routercommon.Domain_RootDomain
		case "full":
			domain.Type = routercommon.Domain_Full
		case "keyword":
			domain.Type = routercommon.Domain_Plain
		case "regexp":
			domain.Type = routercommon.Domain_Regex
		default:
			domain.Type = routercommon.Domain_RootDomain
		}

		for _, attr := range entry.Attrs {
			domain.Attribute = append(domain.Attribute, &routercommon.Domain_Attribute{
				Key: attr,
				TypedValue: &routercommon.Domain_Attribute_BoolValue{
					BoolValue: true,
				},
			})
		}

		site.Domain = append(site.Domain, domain)
	}

	return site
}

func toPlainText(list *ParsedList) string {
	var lines []string
	for _, entry := range list.Entries {
		lines = append(lines, entry.Value)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func main() {
	flag.Parse()

	fmt.Println("Loading domain lists...")
	if err := Load(*dataPath); err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		os.Exit(1)
	}

	geoSiteList := &routercommon.GeoSiteList{}
	var names []string
	for name := range loadedLists {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		parsed, err := ParseList(name, nil)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", name, err)
			continue
		}

		geoSite := toProto(parsed)
		geoSiteList.Entry = append(geoSiteList.Entry, geoSite)
		fmt.Printf("Processed: %s (%d domains)\n", name, len(geoSite.Domain))
	}

	// Sort entries by country code
	sort.Slice(geoSiteList.Entry, func(i, j int) bool {
		return geoSiteList.Entry[i].CountryCode < geoSiteList.Entry[j].CountryCode
	})

	// Marshal to protobuf
	data, err := proto.Marshal(geoSiteList)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		os.Exit(1)
	}

	// Write output
	outputPath := filepath.Join(*outputDir, *outputName)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nGenerated: %s (%d bytes, %d categories)\n", outputPath, len(data), len(geoSiteList.Entry))

	// Export plaintext if requested
	if *exportLists != "" {
		for _, name := range strings.Split(*exportLists, ",") {
			name = strings.TrimSpace(name)
			parsed, err := ParseList(name, nil)
			if err != nil {
				fmt.Printf("Warning: cannot export %s: %v\n", name, err)
				continue
			}

			txtPath := filepath.Join(*outputDir, name+".txt")
			if err := os.WriteFile(txtPath, []byte(toPlainText(parsed)), 0644); err != nil {
				fmt.Printf("Warning: cannot write %s: %v\n", txtPath, err)
				continue
			}
			fmt.Printf("Exported: %s\n", txtPath)
		}
	}
}
