package main

import (
	"flag"
	"log"

	"github.com/ramyad/tucows/internal/api/facade"
	"github.com/ramyad/tucows/internal/app/web"
)

func main() {
	port := flag.Int("port", 8080, "Port number for the web application")
	flag.Parse()

	api := facade.NewAPIFacade()
	app := web.NewWebApp(api, *port)

	err := app.Run()
	if err != nil {
		log.Fatalf("Failed to run terminal application: %v", err)
	}
}
