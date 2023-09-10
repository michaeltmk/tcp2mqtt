package main

import (
	"log"
	"net"
	"time"
)

func send(conn net.Conn) {
	// lets create the message we want to send accross
	// text := `{"IMEI":866250067083620,"adc":[7,15],"addr":"1","cmd":"upload","count":[0,0,0,0,0,0,0,0],"ipFrom":"182.239.107.2","port":"49703", "` + time.Now().Format(time.RFC3339) + `":"123" }`
	text := `{"IMEI":866250067083620, "` + time.Now().Format(time.RFC3339) + `":"123" }`
	log.Printf("msg: %s", text)
	msg := []byte(text)
	conn.Write(msg)
}

func main() {

	for {
		conn, _ := net.Dial("tcp", ":7777")
		send(conn)
		// log.Printf("msg sent")
		time.Sleep(10000 * time.Microsecond)
	}

}
