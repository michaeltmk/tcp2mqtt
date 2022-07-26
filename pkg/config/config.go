package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func GetMqttSchema() ([]byte, error) {
	type V struct {
		MQTT struct {
			SCHEMA struct {
				MESSAGE string
			}
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
	return []byte(v.MQTT.SCHEMA.MESSAGE), nil
}

func GetUsernameSchema() ([]byte, error) {
	type V struct {
		MQTT struct {
			SCHEMA struct {
				USERNAME string
			}
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
	return []byte(v.MQTT.SCHEMA.USERNAME), nil
}

func GetPasswordSchema() ([]byte, error) {
	type V struct {
		MQTT struct {
			SCHEMA struct {
				PASSWORD string
			}
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
	return []byte(v.MQTT.SCHEMA.PASSWORD), nil
}
