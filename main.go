package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/shoarai/washout/jaxfilter"
	"github.com/shoarai/washout/washloop"

	"./internal/receiver"
	"./internal/webserver"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP Address")
	port := flag.Int("port", 8888, "Port number")
	webAddress := flag.String("web-address", ":8080", "Address for Web Server")
	interval := flag.Uint("interval", 10, "Inteval of washout")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *ip, *port)
	err := receiver.Listen(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Listen motion from %q", address))

	done := make(chan struct{})
	loop := createWashloop(*interval)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		fmt.Println()
		fmt.Println("exit...")

		webserver.Close()
		receiver.Close()
		loop.Stop()
		close(done)
	}()

	go func() {
		loop.Start()
	}()

	go func() {
		for {
			data, err := receiver.Read()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(data)

			motion := toMotion(data)
			loop.SetMotion(motion)
		}
	}()

	webserver.ListenAndServe(*webAddress, loop)
	<-done
}

func createWashloop(interval uint) *washloop.WashoutLoop {
	wash := jaxfilter.NewWashout(interval)
	return washloop.NewWashoutLoop(wash, interval)
}

func toMotion(str string) washloop.Motion {
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

	motion := washloop.Motion{
		Acceleration: washloop.Vector{
			X: floats[0],
			Y: floats[1],
			Z: floats[2]},
		AngularVelocity: washloop.Vector{
			X: floats[3],
			Y: floats[4],
			Z: floats[5]},
	}
	return motion
}
