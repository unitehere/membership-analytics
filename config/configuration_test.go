package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
