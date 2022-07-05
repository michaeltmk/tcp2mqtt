package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func GetMqttSchema() ([]byte, error) {
	type V struct {
		VERSION string
		MQTT    struct {
			SCHEMA string
		}
	}
	y, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return nil, err
	}
	v := V{}
	// m := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(y), &v)
	if err != nil {
		return nil, err
	}
	return []byte(v.MQTT.SCHEMA), nil
}
