package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("CONFIG_PATH", "config.yaml")
	result, err := GetMqttSchema()
	assert.NoError(t, err)
	assert.Equal(t, string(`{"{{.netData.IMEI}}":[{"body":{{.netData}},"ipFrom":{{.ip}},"port":{{.port}}}]}`), strings.TrimSpace(string(result)))
}
