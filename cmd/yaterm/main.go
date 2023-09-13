package main

import (
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

func main() {
	// Setup
	app := app.New()
	window := app.NewWindow("YaTerm")

	textGrid := widget.NewTextGrid()
	textGrid.SetText("Hello, World!")

	// Run Process
	command := exec.Command("/bin/zsh")
	process, err := pty.Start(command)

	if err != nil {
		fyne.LogError("Failed to start pty", err)
		os.Exit(1)
	}

	defer command.Process.Kill()

	process.Write([]byte("ls\r"))
	time.Sleep(1 * time.Second)
	b := make([]byte, 1024)
	_, err = process.Read(b)
	if err != nil {
		fyne.LogError("Failed to read pty", err)
	}

	textGrid.SetText(string(b))

	// Display
	window.SetContent(
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(800, 600)),
			textGrid,
		),
	)

	window.ShowAndRun()
}
