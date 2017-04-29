package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	router    *mux.Router
	startTime time.Time

	ready   *sync.Cond
	waiting bool
	geoloc  *GeolocationServer
}

func NewHttpServer() *HttpServer {

	server := new(HttpServer)
	server.startTime = time.Now()

	server.waiting = true
	server.ready = sync.NewCond(&sync.Mutex{})

	server.geoloc = nil

	return server

}

func (server *HttpServer) close() {

	fmt.Println("Good Bye!")
}

func writeResponse(w http.ResponseWriter, body []byte, startTime time.Time) {

	l := len(body)
	w.Header().Set("Server", "Sunny-geoip")
	w.Header().Set("X-Powered-By", "Sunny-geoip/1.0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Length", fmt.Sprint(l))
	w.Header().Set("X-Time", fmt.Sprintf("%s", time.Now().Sub(startTime)))
	w.Write(body)

}

func notFound(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	response := map[string]string{
		"status":  "error",
		"message": "404",
	}
	body, _ := json.Marshal(response)
	writeResponse(w, body, startTime)

}

func (server *HttpServer) bind(g *GeolocationServer) {
	server.geoloc = g
	server.ready.L.Lock()
	server.ready.Broadcast()
	server.waiting = false
	server.ready.L.Unlock()

}

func addCorsHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fn(w, r)
	}
}

func (server *HttpServer) listen(port string) {

	router := mux.NewRouter()
	server.router = router
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/geoip/location", server.geoipHandler)

	go func() {
		fmt.Println("Listening to HTTP port " + port)
		if err := http.ListenAndServe(port, router); err != nil {
			fmt.Printf("Error while listening: %s\n", err)
			os.Exit(1)
		}

	}()
}
