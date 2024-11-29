package debug

import (
	"fmt"
	"net"
)

type debugConn struct {
	lastMessageIndex int
	conn             net.Conn
}

type debugServer struct {
	messages    []string
	listener    net.Listener
	connections []debugConn
}

var ds debugServer

var messages chan string = make(chan string, 100)
var quit chan bool = make(chan bool)

func Start() error {
	var err error
	ds.listener, err = net.Listen("tcp", "localhost:13327")
	if err != nil {
		println("Failed to start debug server")
	}
	go func() {
		for {
			conn, err := ds.listener.Accept()
			if err != nil {
				println("Failed to accept debug connection")
			}
			ds.connections = append(ds.connections, debugConn{0, conn})
			for _, msg := range ds.messages {
				conn.Write([]byte(msg + "\n"))
			}
		}
	}()
	go func() {
		for {
			select {
			case msg := <-messages:
				ds.messages = append(ds.messages, msg)
				for _, conn := range ds.connections {
					conn.conn.Write([]byte(msg + "\n"))
				}
			case <-quit:
				return
			}
		}
	}()
	return nil
}

func Stop() {
	if err := ds.listener.Close(); err == nil {
		<-quit
		return
	}
}

func Debug(v ...interface{}) {
	messages <- fmt.Sprint(v...)
}

func BytesToStringAndHex(b []byte) string {
	result := ""
	for _, c := range b {
		result += fmt.Sprintf("%02x ", c)
	}
	return result
}
