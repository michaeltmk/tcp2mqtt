# TCP2MQTT

A proxy to recieve raw TCP sockets and send to a MQTT broker with customable format.
Froked from [tcp2mqtt](https://github.com/gonzalo123/tcp2mqtt)

It is a go client that reads the TCP sockets and send the information to the MQTT broker.
It support json format as TCP sockets foramt only.

## Getting Started
1. Edit the coustomated schema in config.yaml file
	```yaml
	version: 1
	mqtt:
	schema:
		message: |
		{{. | fjson}}
		username: |
		{{.IMEI | printf "%.f" -}}
		password: ""
	```
2. Enter the MQTT broker configuration in environment.BROKER in docker-compose.yml
```yaml
environment:
	- CONFIG_PATH=/opt/config.yaml
	- BROKER=tcp://localhost:1883
```
run ```docker-compose up --build```



## go-template
We use go-template to generate the MQTT message.
There is a customated function imported into the template engine.

```go
// orderedMarshalString marshals a given value into JSON with ordered keys
func orderedMarshalString(v any) (string, error) {
	b, err := encoder.Encode(v, encoder.SortMapKeys)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
```
``` go
template.FuncMap{
	"fjson": orderedMarshalString,
}
```
