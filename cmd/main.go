package main

import (
	"github.com/pasDamola/file-tracker/config"
	"github.com/pasDamola/file-tracker/daemon"
	"github.com/pasDamola/file-tracker/logs"
	"github.com/pasDamola/file-tracker/ui"
)

func main() {
	logs.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logs.Log(err.Error())
		return
	}

	d := daemon.NewDaemon(cfg)

	// go api.StartAPI(d)

	ui.StartUI(d) // Start the UI.
}
