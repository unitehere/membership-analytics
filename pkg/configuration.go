package config

import (
	"github.com/tkanos/gonfig"
)

// The Configuration type describes the configuration variables
type Configuration struct {
	ElasticUsername string
	ElasticPassword string
}

// Values is the access point to the configuration values
var Values = Configuration{}

// Init loads the config
func Init() {
	err := gonfig.GetConf("config.json", &Values)
	if err != nil {
		panic("Could not load configuration file. Please check that you have a valid config.json")
	}
}
