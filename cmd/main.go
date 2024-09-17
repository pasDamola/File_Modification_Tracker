package main

import (
	"log"

	"github.com/pasDamola/file-tracker/config"
	"github.com/pasDamola/file-tracker/internal/adapters/http"
	"github.com/pasDamola/file-tracker/internal/adapters/osquery"
	"github.com/pasDamola/file-tracker/internal/adapters/ui"
	"github.com/pasDamola/file-tracker/internal/core/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	fileService := services.NewFileService()
	osqueryAdapter, err := osquery.NewOsqueryAdapter(cfg.SocketPath)
	if err != nil {
		log.Fatalf("Error creating osquery adapter: %v", err)
	}

	d := services.NewDaemon(cfg, osqueryAdapter, fileService)

	go http.StartAPI(fileService, d)

	uiApp := ui.NewApp(d)
	uiApp.Start()
}
