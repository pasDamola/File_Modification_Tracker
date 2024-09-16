package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pasDamola/file-tracker/api"
	"github.com/pasDamola/file-tracker/daemon"
)

func StartUI(d *daemon.Daemon) {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Modification Tracker")

	startButton := widget.NewButton("Start Service", func() {
		go func() {
			d.Start() // Start the daemon
			// Start the API in a separate goroutine
			go api.StartAPI(d) // Assuming you have a StartAPI method in your Daemon struct
			widget.ShowPopUp(widget.NewLabel("Service Started"), myWindow.Canvas())
		}()
	})

	stopButton := widget.NewButton("Stop Service", func() {
		close(d.WorkerQueue)
		widget.ShowPopUp(widget.NewLabel("Service Stopped"), myWindow.Canvas())
	})

	myWindow.SetContent(container.NewVBox(
		startButton,
		stopButton,
	))

	size := fyne.NewSize(300, 200)
	myWindow.Resize(size)

	myWindow.ShowAndRun()
}
