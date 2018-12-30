package webserver

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"text/template"

	"github.com/shoarai/washout/washloop"

	"../../models"
	"github.com/shoarai/washout"
)

var loop *washloop.WashoutLoop
var listener net.Listener

// ListenAndServe listens and serve requests as web server.
func ListenAndServe(address string, l *washloop.WashoutLoop) error {
	loop = l
	return listenAndServe(address)
}

func listenAndServe(address string) error {
	var err error
	listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/position", position)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	return http.Serve(listener, mux)
}

// Close closes a listener of server
func Close() {
	if listener != nil {
		listener.Close()
	}
}

func startWebServer() {
	http.HandleFunc("/", index)
	http.HandleFunc("/position", position)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
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

func position(w http.ResponseWriter, r *http.Request) {
	position := loop.GetPosition()
	response, err := json.Marshal(toPosition(position))
	if err != nil {
		panic(err)
	}
	// fmt.Fprintf(w, response)
	w.Write(response)
}

func toPosition(pos washout.Position) models.Position {
	return models.Position{
		X:      pos.X,
		Y:      pos.Y,
		Z:      pos.Z,
		AngleX: pos.AngleX,
		AngleY: pos.AngleY,
		AngleZ: pos.AngleZ,
	}
}
