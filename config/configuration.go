package config

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/tkanos/gonfig"
)

// The Configuration type describes the configuration variables
type Configuration struct {
	ElasticUsername string
	ElasticPassword string
}

// Values is the access point to the configuration values
var Values = Configuration{}

func init() {
	err := gonfig.GetConf(getFileName(), &Values)
	if err != nil {
		panic(err)
	}
}

// must use an absolute path due to testing. Will cause problems if only "config.json"
func getFileName() string {
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), "config.json")
	return filePath
}
