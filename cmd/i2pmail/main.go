package main

import (
	"flag"
	"log"
	"net"
	"net/http"
)

var (
	port      = flag.String("port", "8080", "port to listen on")
	host      = flag.String("host", "i2pmail.org", "host to listen on")
	directory = flag.String("directory", "./www", "directory to serve")
)

func main() {
	fs := http.FileServer(http.Dir(*directory))
	//http.Handle("/", fs)

	address := net.JoinHostPort(*host, *port)
	log.Printf("Listening on %s...", address)
	log.Printf("Serving %s...", *directory)
	log.Printf("Args were %s, %s, %s", *port, *host, *directory)
	err := http.ListenAndServe(address, fs)
	if err != nil {
		log.Fatal(err)
	}
}
