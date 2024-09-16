package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pasDamola/file-tracker/api"
	"github.com/pasDamola/file-tracker/daemon"
)

func StartUI(d *daemon.Daemon) {
	fmt.Println("Starting Fyne UI")
	myApp := app.New()
	myWindow := myApp.NewWindow("File Modification Tracker")

	startButton := widget.NewButton("Start Service", func() {
		go func() {
			d.Start()
			go api.StartAPI(d)
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
	fmt.Println("Where we are")

	size := fyne.NewSize(300, 200)
	myWindow.Resize(size)

	myWindow.ShowAndRun()
}
