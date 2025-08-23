package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/melsincostan/menhir/menhir"
	"github.com/melsincostan/menhir/modules/cors"
	"github.com/melsincostan/menhir/modules/logging"
	"github.com/melsincostan/menhir/modules/xffor"
)

func main() {
	destination := flag.String("destination", "http://localhost", "reverse proxy target")
	host := flag.String("host", "0.0.0.0", "host on which the reverse proxy will listen")
	port := flag.String("port", "8080", "port on which the reverse proxy will listen")
	listModules := flag.Bool("list-modules", false, "list all registered modules")

	wrapper := menhir.New()
	wrapper.Register(cors.New(), &logging.Logging{}, xffor.New())

	modEnableArgs := map[string]*bool{}

	for _, mod := range wrapper.Modules() {
		modEnableArgs[mod.Name()] = flag.Bool(fmt.Sprintf("module.%s", mod.Name()), mod.Default(), fmt.Sprintf("activate the %s module", mod.Name()))
	}

	flag.Parse()

	if *listModules {
		fmt.Print("Registered modules:\n")
		for _, mod := range wrapper.Modules() {
			fmt.Printf("- %s\n", mod.Name())
		}
		os.Exit(0)
	}

	for mname, enabled := range modEnableArgs {
		if *enabled {
			if err := wrapper.Enable(mname); err != nil {
				log.Fatalf("could not enable modules: %s", err.Error())
			}
		}
	}

	if err := wrapper.Init(*destination); err != nil {
		log.Fatalf("Could not start: %s", err.Error())
	}

	http.ListenAndServe(net.JoinHostPort(*host, *port), wrapper)
}
