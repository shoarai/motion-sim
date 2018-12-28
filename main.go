package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/shoarai/washout"

	"./receiver"

	"github.com/alecthomas/template"
	"github.com/shoarai/washout/jaxfilter"
)

type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	AngleX float64 `json:"angleX"`
	AngleY float64 `json:"angleY"`
	AngleZ float64 `json:"angleZ"`
}

var motionReceiver = receiver.MotionReceiver{}

func main() {
	// ip := "127.0.0.1"
	ip := "192.168.179.3"
	port := 8885

	address := ip + ":" + strconv.Itoa(port)
	done := make(chan struct{})

	err, errCh := motionReceiver.Listen(address, done)
	if err != nil {
		log.Fatal(err)
	}

	listener, ch, err := ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// os.Stdin.Read(make([]byte, 1)) // read a single byte
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		fmt.Println()
		fmt.Println("exit...")
		listener.Close()
		close(done)
	}()

	<-ch
	<-errCh
	// log.Fatal(<-ch)
	// log.Fatal(<-cherr)
}

func ListenAndServe() (net.Listener, chan error, error) {
	ch := make(chan error)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return listener, nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/position", position)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	go func() {
		ch <- http.Serve(listener, mux)
	}()
	return listener, ch, nil
}

func startWebServer() {
	http.HandleFunc("/", viewHandler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/position", position)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}

func position(w http.ResponseWriter, r *http.Request) {
	motion := motionReceiver.GetMotion()

	filter := jaxfilter.NewWashout(100)
	position := filter.Filter(
		motion.Acceleration.X,
		motion.Acceleration.Y,
		motion.Acceleration.Z,
		motion.AngularVelocity.X,
		motion.AngularVelocity.Y,
		motion.AngularVelocity.Z,
	)

	response, err := json.Marshal(toPosition(position))
	if err != nil {
		panic(err)
	}
	// fmt.Fprintf(w, response)
	w.Write(response)
}

func toPosition(pos washout.Position) Position {
	return Position{
		pos.X, pos.Y, pos.Z, pos.AngleX, pos.AngleY, pos.AngleZ,
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// page := Page{"Hello World.", 1}
	tmpl, err := template.ParseFiles("index.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
	// err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}
