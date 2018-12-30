package receiver

import (
	"net"
)

var connection net.PacketConn

func Listen(address string) (err error) {
	connection, err = net.ListenPacket("udp", address)
	return
}

func Read() (string, error) {
	buffer := make([]byte, 1500)
	length, _, err := connection.ReadFrom(buffer)
	if err != nil {
		return "", err
	}
	// fmt.Println(string(buffer[:length]))
	// fmt.Printf("Received from %v: %v\n", remoteAddr, string(buffer[:length]))
	// conn.WriteTo([]byte("Hello, World !"), remoteAddr)

	return string(buffer[:length]), nil
}

func Close() {
	connection.Close()
}
