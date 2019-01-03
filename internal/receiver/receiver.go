package receiver

import (
	"net"
)

var connection net.PacketConn

// Listen listens an UDP connection.
func Listen(address string) (err error) {
	connection, err = net.ListenPacket("udp", address)
	return
}

// Read reads data of a connection.
func Read() (string, error) {
	buffer := make([]byte, 1500)
	length, _, err := connection.ReadFrom(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:length]), nil
}

// Close closes a connection.
func Close() {
	connection.Close()
}
