package main

import (
	"log"
	"net"
	"time"
)

func send(conn net.Conn) {
	// lets create the message we want to send accross
	msg := []byte(`{"IMEI":866250067083620,"adc":[7,15],"addr":"1","cmd":"upload","count":[0,0,0,0,0,0,0,0],"ipFrom":"182.239.107.2","port":"49703"}`)
	conn.Write(msg)
}

func main() {

	for {
		conn, _ := net.Dial("tcp", ":7777")
		send(conn)
		log.Printf("msg sent")
		time.Sleep(100000 * time.Microsecond)
	}

}
