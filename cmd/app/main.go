package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"tcp2mqtt/pkg/config"
	"tcp2mqtt/pkg/message"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	port, closeConnection, topic, broker := parseFlags()
	openSocket(*port, *closeConnection, *topic, *broker, onMessage)
}

func openSocket(port string, closeConnection bool, topic string, broker string, onMessage func(url string, username string, password string, topic string, buffer string)) {
	msgSchema, err := config.GetMqttSchema()
	if err != nil {
		log.Fatalln(err)
	}
	usernameSchema, err := config.GetMqttSchema()
	if err != nil {
		log.Fatalln(err)
	}
	passwordSchema, err := config.GetMqttSchema()
	if err != nil {
		log.Fatalln(err)
	}
	PORT := "0.0.0.0:" + port
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
		go handleConnection(c, msgSchema, usernameSchema, passwordSchema, closeConnection, topic, broker, onMessage)
	}
}

func createClientOptions(url string, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetUsername(username)
	opts.SetPassword(password)
	return opts
}

func connect(url string, username string, password string) mqtt.Client {
	opts := createClientOptions(url, username, password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func onMessage(url string, username string, password string, topic string, buffer string) {
	client := connect(url, username, password)
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

func handleConnection(c net.Conn, msgSchema []byte, usernameSchema []byte, passwordSchema []byte, closeConnection bool, topic string, broker string, onMessage func(url string, username string, password string, topic string, buffer string)) {
	log.Printf("Accepted connection from %s\n", c.RemoteAddr().String())
	for {
		ip, port, err := net.SplitHostPort(c.RemoteAddr().String())
		if err != nil {
			log.Println(err)
		}
		netData, err := bufio.NewReader(c).ReadString('}')
		if err != nil {
			log.Printf("error reading netData: %v", err)
		}
		log.Printf("netData: %v", netData)
		handledMessage, err := message.HandleWholeMsg(ip, port, netData, msgSchema, usernameSchema, passwordSchema)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("sending to topic %s message:%s\n", topic, handledMessage.Msg)
			onMessage(broker, handledMessage.Username, handledMessage.Password, topic, handledMessage.Msg)
			log.Printf("sent")
		}
		if closeConnection {
			c.Close()
			return
		}
	}
}
