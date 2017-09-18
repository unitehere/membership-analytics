package config

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

// The Configuration type describes the configuration variables
type Configuration struct {
	ElasticUsername string
	ElasticPassword string
	Index           string
}

// Values is the access point to the configuration values
var Values = Configuration{}

func init() {
	err := gonfig.GetConf(getFilePath(), &Values)
	if err != nil {
		panic(err)
	}
}

// must use an absolute path due to testing. Will cause problems if only "config.json"
func getFilePath() string {
	envName := "dev"
	if strings.ToLower(os.Getenv("ENV")) == "prod" || strings.ToLower(os.Getenv("ENV")) == "production" {
		envName = "prod"
	}
	fileName := "config." + envName + ".json"

	_, dirname, _, _ := runtime.Caller(1)
	filePath := path.Join(path.Dir(dirname), "..", fileName)
	return filePath
}
