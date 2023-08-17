package main

import (
	"log"

	"github.com/ramyad/tucows/internal/api/facade"
	"github.com/ramyad/tucows/internal/app/terminal"
)

func main() {
	api := facade.NewAPIFacade()
	app := terminal.NewTerminalApp(api)
	err := app.Run()
	if err != nil {
		log.Fatalf("Failed to run terminal application: %v", err)
	}
}
