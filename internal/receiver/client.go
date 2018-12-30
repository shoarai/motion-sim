package receiver

import (
	"fmt"
	"net"
)

func Send() {
	conn, _ := net.Dial("udp", "192.168.179.3:8888")
	// conn, _ := net.Dial("udp", "127.0.0.1:8888")
	defer conn.Close()
	fmt.Println("サーバへメッセージを送信.")
	conn.Write([]byte("Hello From Client."))

	fmt.Println("サーバからメッセージを受信。")
	buffer := make([]byte, 1500)
	length, _ := conn.Read(buffer)
	fmt.Printf("Receive: %s \n", string(buffer[:length]))
}
