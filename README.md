<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://varadise.ai/">
    <img src="docs/images/logo.webp" alt="Logo" />
  </a>
  <h3 align="center">TCP2MQTT</h3>
  <p align="center">
    <br />

  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="LICENSE.txt">License</a></li>
    <li><a href="#maintainers">Maintainers</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

A proxy to recieve raw TCP sockets and send to a MQTT broker with customable format.
Froked from [tcp2mqtt](https://github.com/gonzalo123/tcp2mqtt)

It is a go client that reads the TCP sockets and send the information to the MQTT broker.
It support json format as TCP sockets foramt only.


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
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

### Prerequisites


### Installation


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MAINTAINERS -->
## Maintainers

An up-to-date list of people involved with development / support of the project

<!-- Simon Ball - [simonball@varadise.cloud](mailto:simonball@varadise.cloud) -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[product-screenshot]: docs/images/screenshot.png
<!-- Links to pages for products / components commonly used across the company -->
<!-- If you're project is using something not listed below, please create a ticket -->
<!-- in the Varadise Kitchen Board to propose updating the master template -->
[Next-url]: https://nextjs.org/
[React-url]: https://reactjs.org/
[Vue-url]: https://vuejs.org/
[Flask-url]: https://flask.palletsprojects.com
[Gin-url]: https://gin-gonic.com/