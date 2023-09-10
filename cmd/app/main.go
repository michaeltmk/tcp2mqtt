package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"runtime"
	"tcp2mqtt/pkg/config"
	"tcp2mqtt/pkg/message"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
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
	usernameSchema, err := config.GetUsernameSchema()
	if err != nil {
		log.Fatalln(err)
	}
	passwordSchema, err := config.GetPasswordSchema()
	if err != nil {
		log.Fatalln(err)
	}
	messagetypeSchema, err := config.GetMessageTypeSchema()
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
			log.Println(err)
			return
		}
		go handleConnection(c, msgSchema, usernameSchema, passwordSchema, messagetypeSchema, closeConnection, topic, broker, onMessage)
	}
}

func createClientOptions(url string, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(uuid.New().String())
	// opts.SetKeepAlive(2* time.Minute)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(500 * time.Millisecond)
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
	t := client.Publish(topic, 1, true, buffer)
	go func() {
		_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if t.Error() != nil {
			log.Printf("error in sending mqtt data: %v", t.Error()) // Use your preferred logging technique (or just fmt.Printf)
		}
	}()
	log.Println(t)
	defer client.Disconnect(250)
}

func parseFlags() (*string, *bool, *string, *string) {
	port := flag.String("port", "7777", "port number")
	closeConnection := flag.Bool("close", true, "Close connection")
	topic := flag.String("topic", "topic", "mqtt topic")
	broker := flag.String("broker", "tcp://localhost:1883", "mqtt server")
	flag.Parse()

	return port, closeConnection, topic, broker
}

func handleConnection(c net.Conn, msgSchema []byte, usernameSchema []byte, passwordSchema []byte, messagetypeSchema []byte, closeConnection bool, topic string, broker string, onMessage func(url string, username string, password string, topic string, buffer string)) {
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
		handledMessage, err := message.HandleWholeMsg(ip, port, netData, msgSchema, usernameSchema, passwordSchema, messagetypeSchema)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("sending to topic %s message:%s\n", topic, handledMessage.Msg)
			onMessage(broker, handledMessage.Username, handledMessage.Password, topic, handledMessage.Msg)
			log.Printf("sent")
		}
		if closeConnection {
			log.Printf("Closing connection from %s\n", c.RemoteAddr().String())
			c.Close()
			runtime.GC()
			return
		}
	}
}
