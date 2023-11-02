package main

import (
	"flag"

	"esd/internal"
	"esd/resolver"
)

var (
	configFile = "conf/esd.yml"
)

func main() {
	flag.StringVar(&configFile, "config", configFile, "path to YAML config")
	config, err := internal.NewConfig(configFile)
	if err != nil {
		panic(err)
	}
	defaultZone := &resolver.Zone{
		Logger:  config.Logger,
		Root:    ".",
		Parent:  config.Parent,
		Records: nil,
	}
	config.Zones = append(config.Zones, defaultZone)
	config.Logger.Info(resolver.Start(config.Listen, config.Protocol, config.Zones))
}
