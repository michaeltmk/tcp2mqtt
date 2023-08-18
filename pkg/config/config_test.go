package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMqttSchema(t *testing.T) {
	os.Setenv("CONFIG_PATH", "config.yaml")
	result, err := GetMqttSchema()
	assert.NoError(t, err)
	assert.Equal(t, string(`{"rawMsg":{{. | fjson}},"ipFrom":"{{.ip}}","port":"{{.port}}"}`), strings.TrimSpace(string(result)))
}

func TestGetUsernameSchema(t *testing.T) {
	os.Setenv("CONFIG_PATH", "config.yaml")
	result, err := GetUsernameSchema()
	assert.NoError(t, err)
	assert.Equal(t, string(`{{.IMEI | printf "%.f"}}`), strings.TrimSpace(string(result)))
}

func TestGetPasswordSchema(t *testing.T) {
	os.Setenv("CONFIG_PATH", "config.yaml")
	result, err := GetPasswordSchema()
	assert.NoError(t, err)
	assert.Equal(t, string(``), strings.TrimSpace(string(result)))
}

func TestGetMessageTypeSchema(t *testing.T) {
	os.Setenv("CONFIG_PATH", "config.yaml")
	result, err := GetMessageTypeSchema()
	assert.NoError(t, err)
	assert.Equal(t, string(`json`), strings.TrimSpace(string(result)))
}
