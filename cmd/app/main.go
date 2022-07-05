package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"tcp2mqtt/pkg/config"
	"tcp2mqtt/pkg/template"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	port, closeConnection, topic, broker := parseFlags()
	openSocket(*port, *closeConnection, *topic, *broker, onMessage)
}

func openSocket(port string, closeConnection bool, topic string, broker string, onMessage func(url string, topic string, buffer string)) {
	s, err := config.GetMqttSchema()
	if err != nil {
		log.Fatalln(err)
	}
	PORT := "localhost:" + port
	l, err := net.Listen("tcp4", PORT)
	log.Printf("Serving %s\n", l.Addr().String())
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConnection(c, s, closeConnection, topic, broker, onMessage)
	}
}

func createClientOptions(url string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	return opts
}

func connect(url string) mqtt.Client {
	opts := createClientOptions(url)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func onMessage(url string, topic string, buffer string) {
	client := connect(url)
	token := client.Publish(topic, 2, false, buffer)
	log.Println(token)
}

func parseFlags() (*string, *bool, *string, *string) {
	port := flag.String("port", "7777", "port number")
	closeConnection := flag.Bool("close", true, "Close connection")
	topic := flag.String("topic", "topic", "mqtt topic")
	broker := flag.String("broker", "tcp://localhost:1883", "mqtt topic")
	flag.Parse()

	return port, closeConnection, topic, broker
}

func handleMsg(ip string, port string, netData string, schema []byte) ([]byte, error) {
	netDataJS := make(map[string]any)
	err := json.Unmarshal([]byte(strings.TrimSpace(netData)), &netDataJS)
	if err != nil {
		return nil, err
	}
	msgData := map[string]interface{}{
		"netData": netDataJS,
		"ip":      ip,
		"port":    port,
	}
	bytesRepresentation, err := json.Marshal(msgData)
	if err != nil {
		return nil, err
	}
	result, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handleConnection(c net.Conn, schema []byte, closeConnection bool, topic string, broker string, onMessage func(url string, topic string, buffer string)) {
	log.Printf("Accepted connection from %s\n", c.RemoteAddr().String())
	for {
		ip, port, err := net.SplitHostPort(c.RemoteAddr().String())
		if err != nil {
			log.Println(err)
		}
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("error reading netData: %v", err)
		}
		result, err := handleMsg(ip, port, netData, schema)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("sending to topic %s message:%s\n", topic, result)
			onMessage(broker, topic, string(result))
			log.Printf("sent")
		}
		if closeConnection {
			c.Close()
			return
		}
	}
}
