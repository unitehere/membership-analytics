package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	assert.NotEmpty(t, Values.ElasticUsername)
	assert.NotEmpty(t, Values.ElasticPassword)
	assert.NotEmpty(t, Values.Index)
}

func TestGetFilePath(t *testing.T) {
	defer os.Unsetenv("ENV")
	assert.Contains(t, getFilePath(), "github.com/unitehere/membership-analytics/config/config.dev.json")

	os.Setenv("ENV", "prod")
	assert.Contains(t, getFilePath(), "github.com/unitehere/membership-analytics/config/config.prod.json")
}
