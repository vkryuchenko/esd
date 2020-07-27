package internal

import (
	"esd/resolver"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type AppConfig struct {
	Listen    string          `yaml:"listen"`
	Protocol  string          `yaml:"protocol"`
	ParentDNS []string        `yaml:"parentDns"`
	Zones     []resolver.Zone `yaml:"zones"`
}

func (ac *AppConfig) Read(filePath string) error {
	configPath, err := ResolvePath(filePath)
	if err != nil {
		return err
	}
	cf, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer cf.Close()
	err = yaml.NewDecoder(cf).Decode(ac)
	if err != nil {
		return err
	}
	ac.Zones = ac.normalizeZones(ac.Zones)
	return err
}

func (ac *AppConfig) normalizeZones(zones []resolver.Zone) []resolver.Zone {
	var normalizedZones []resolver.Zone
	for _, zone := range zones {
		nZ := resolver.Zone{
			Root:      zone.Root,
			ParentDNS: zone.ParentDNS,
			Records:   []resolver.Record{},
		}
		if !strings.HasSuffix(zone.Root, ".") {
			nZ.Root = zone.Root + "."
		}
		for _, record := range zone.Records {
			nR := resolver.Record{
				Name:    strings.TrimSuffix(record.Name, ".") + "." + nZ.Root,
				Type:    strings.ToUpper(record.Type),
				PointTo: record.PointTo,
			}
			if record.Name == zone.Root {
				nR.Name = nZ.Root
			}
			nZ.Records = append(nZ.Records, nR)
		}
		if len(nZ.ParentDNS) < 1 {
			nZ.ParentDNS = ac.ParentDNS
		}
		normalizedZones = append(normalizedZones, nZ)
	}
	return normalizedZones
}
