package internal

import (
	"log"
	"os"
	"strings"

	"esd/resolver"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger    *zap.SugaredLogger `yaml:"-"`
	LogLevel  string             `yaml:"logLevel"`
	Listen    string             `yaml:"listen"`
	Protocol  string             `yaml:"protocol"`
	Parent    resolver.Parent    `yaml:"parent"`
	ZoneFiles []string           `yaml:"zoneFiles"`
	Zones     []*resolver.Zone   `yaml:"zones"`
}

func NewConfig(filePath string) (*Config, error) {
	configPath, err := resolvePath(filePath)
	if err != nil {
		return nil, err
	}
	cf, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer cf.Close()
	config := &Config{}
	err = yaml.NewDecoder(cf).Decode(config)
	if err != nil {
		return nil, err
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	config.Logger = NewZapLogger(config.LogLevel)
	for _, zoneFilePath := range config.ZoneFiles {
		zone, err := readZoneFromFile(zoneFilePath)
		if err != nil {
			log.Print(err)
			continue
		}
		config.Zones = append(config.Zones, zone)
	}
	config.normalizeZones()
	return config, err
}

func readZoneFromFile(filePath string) (*resolver.Zone, error) {
	configPath, err := resolvePath(filePath)
	if err != nil {
		return nil, err
	}
	zoneFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer zoneFile.Close()
	zone := &resolver.Zone{}
	err = yaml.NewDecoder(zoneFile).Decode(zone)
	return zone, err
}

func (c *Config) normalizeZones() {
	var normalizedZones []*resolver.Zone
	for _, zone := range c.Zones {
		nZ := &resolver.Zone{
			Logger:  c.Logger,
			Root:    zone.Root,
			Parent:  zone.Parent,
			Records: []resolver.Record{},
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
		if len(nZ.Parent.Servers) < 1 {
			nZ.Parent = c.Parent
		}
		normalizedZones = append(normalizedZones, nZ)
	}
	c.Zones = normalizedZones
}
