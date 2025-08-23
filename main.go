package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/melsincostan/menhir/menhir"
	"github.com/melsincostan/menhir/modules/cors"
)

func main() {
	corsModule := cors.New()
	destination := flag.String("destination", "localhost", "reverse proxy target")
	host := flag.String("host", "0.0.0.0", "host on which the reverse proxy will listen")
	port := flag.String("port", "8080", "port on which the reverse proxy will listen")
	wrapper := menhir.New()

	wrapper.Register(corsModule)

	flag.Parse()

	wrapper.Enable("cors")

	if err := wrapper.Init(*destination); err != nil {
		log.Fatalf("Could not start: %s", err.Error())
	}

	http.ListenAndServe(net.JoinHostPort(*host, *port), wrapper)
}
