package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/melsincostan/menhir/menhir"
)

func main() {
	destination := flag.String("destination", "localhost", "reverse proxy target")
	host := flag.String("host", "0.0.0.0", "host on which the reverse proxy will listen")
	port := flag.String("port", "8080", "port on which the reverse proxy will listen")
	wrapper := menhir.New()

	flag.Parse()

	if err := wrapper.Init(*destination); err != nil {
		log.Fatalf("Could not start: %s", err.Error())
	}

	http.ListenAndServe(net.JoinHostPort(*host, *port), wrapper)
}
