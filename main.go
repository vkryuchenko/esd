package main

import (
	"esd/internal"
	"esd/resolver"
	"flag"
	"log"
)

var (
	configFile = "conf/esd.yml"
)

func main() {
	flag.StringVar(&configFile, "config", configFile, "path to YAML config")
	appConfig := internal.AppConfig{}
	err := appConfig.Read(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defaultZone := resolver.Zone{
		Root:    ".",
		Parent:  appConfig.Parent,
		Records: nil,
	}
	appConfig.Zones = append(appConfig.Zones, defaultZone)
	log.Println(resolver.Start(appConfig.Listen, appConfig.Protocol, appConfig.Zones))
}
