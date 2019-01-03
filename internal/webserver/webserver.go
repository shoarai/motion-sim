package webserver

import (
	"encoding/json"
	"net"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shoarai/washout"
	"github.com/shoarai/washout/washloop"
)

// A Position is six degrees of freedom position.
type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	AngleX float64 `json:"angleX"`
	AngleY float64 `json:"angleY"`
	AngleZ float64 `json:"angleZ"`
}

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
	// mux.HandleFunc("/position", position)
	mux.HandleFunc("/ws", wsHandler)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	return http.Serve(listener, mux)
}

// Close closes a listener of server
func Close() {
	if listener != nil {
		listener.Close()
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		panic(err)
	}

	for {
		writePosition(conn)
		time.Sleep(100 * time.Millisecond)
	}
}

func writePosition(conn *websocket.Conn) {
	position := loop.GetPosition()
	pos := toPosition(position)
	conn.WriteJSON(pos)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./index.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
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

func toPosition(pos washout.Position) Position {
	return Position{
		X:      pos.X,
		Y:      pos.Y,
		Z:      pos.Z,
		AngleX: pos.AngleX,
		AngleY: pos.AngleY,
		AngleZ: pos.AngleZ,
	}
}
