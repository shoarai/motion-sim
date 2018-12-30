package receiver

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Vector struct {
	X, Y, Z float64
}

type Motion struct {
	Acceleration    Vector
	AngularVelocity Vector
}

type MotionReceiver struct {
	connection   net.PacketConn
	motionString string
	mutex        *sync.Mutex
}

func (m *MotionReceiver) Listen(address string, done chan struct{}) (error, chan error) {
	m.motionString = "0,0,0,0,0,0"
	return m.receive(address, done)
}

func (m *MotionReceiver) Close() {
	m.connection.Close()
}

func (m *MotionReceiver) GetMotion() *Motion {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	motion := toMotion(m.motionString)
	return motion
}

func toMotion(str string) *Motion {
	strings := strings.Split(str, ",")
	if len(strings) != 6 {
		panic("Invalid length of motion data")
	}
	var floats [6]float64
	for i, v := range strings {
		var err error
		floats[i], err = strconv.ParseFloat(v, 64)
		if err != nil {
			panic("Invalid motion data to convert string to float64")
		}
	}

	motion := Motion{
		Vector{
			floats[0],
			floats[1],
			floats[2]},
		Vector{
			floats[3],
			floats[4],
			floats[5]},
	}
	return &motion
}

func (m *MotionReceiver) receive(address string, done chan struct{}) (error, chan error) {
	fmt.Println("Server is Running at " + address)
	var err error
	m.connection, err = net.ListenPacket("udp", address)
	if err != nil {
		return err, nil
	}

	buffer := make([]byte, 1500)

	m.mutex = new(sync.Mutex)
	ch := make(chan error)

	go func() {
		for {
			length, remoteAddr, err := m.connection.ReadFrom(buffer)
			if err != nil {
				// TODO: Ignore "use of closed network connection"
				ch <- err
				return
			}
			// fmt.Println(string(buffer[:length]))
			fmt.Printf("Received from %v: %v\n", remoteAddr, string(buffer[:length]))
			// conn.WriteTo([]byte("Hello, World !"), remoteAddr)

			m.motionString = string(buffer[:length])
			m.setMotionString(string(buffer[:length]))
		}
	}()

	return nil, ch
}

func (m *MotionReceiver) setMotionString(str string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.motionString = str
}
