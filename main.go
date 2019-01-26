package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
	"log"
	"net/http"
	"os"
	"time"
)

var pc requestInterface = &requestClient{}

func main() {
	// Initialise the pusher connection
	pc.InitConn()

	// Start the router
	r := mux.NewRouter()
	r.HandleFunc("/pour", handlePintRequest)
	_ = http.ListenAndServe(":8000", r)
}

func handlePintRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Beer requested at", time.Now())

	// Allow cross origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")

	pc.SendPour()
	w.WriteHeader(200)
}

type requestInterface interface {
	InitConn()
	SendPour()
}

type requestClient struct {
	client pusher.Client
}

func (pc *requestClient) InitConn() {
	// Get the secret key
	secretKey := os.Getenv("PUSHER_SECRET_KEY")
	if secretKey == "" && flag.Lookup("test.v") == nil {
		log.Fatal("PUSHER_SECRET_KEY not set.")
	}

	client := pusher.Client{
		AppId:   "700326",
		Key:     "fa06261efc5349a70ad5",
		Secret:  secretKey,
		Cluster: "eu",
	}

	pc.client = client
}

func (pc *requestClient) SendPour() {
	_, err := pc.client.Trigger("pour", "pour", nil)
	if err != nil {
		log.Println("Got an error sending item to Pusher channel")
		log.Println(err.Error())
	}
}
