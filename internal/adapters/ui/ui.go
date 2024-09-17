package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pasDamola/file-tracker/internal/core/services"
)

type App struct {
	app    fyne.App
	window fyne.Window
	daemon *services.Daemon
}

func NewApp(d *services.Daemon) *App {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Modification Tracker")
	return &App{
		app:    myApp,
		window: myWindow,
		daemon: d,
	}
}

func (a *App) Start() {
	startButton := widget.NewButton("Start Service", func() {
		a.daemon.Start()
	})

	stopButton := widget.NewButton("Stop Service", func() {
		a.daemon.Stop()
	})

	a.window.SetContent(container.NewVBox(
		startButton,
		stopButton,
	))

	a.window.Resize(fyne.NewSize(300, 200))

	// Close event handler to keep daemon running
	a.window.SetOnClosed(func() {
		dialog.NewConfirm("Confirm", "Are you sure you want to close the application? The service will continue running.", func(confirm bool) {
			if confirm {
				log.Println("UI closed, but the daemon is still running.")
			}

		}, a.window).Show()
	})

	a.window.ShowAndRun()
}
