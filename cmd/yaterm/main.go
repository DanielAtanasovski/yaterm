package main

import (
	"bufio"
	"io"
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

const MaxBufferSize = 16

func main() {
	// Setup
	app := app.New()
	window := app.NewWindow("YaTerm")

	textGrid := widget.NewTextGrid()
	textGrid.SetText("Hello, World!")

	// Run Process
	command := exec.Command("/bin/sh")
	process, err := pty.Start(command)

	if err != nil {
		fyne.LogError("Failed to start pty", err)
		os.Exit(1)
	}

	defer command.Process.Kill()

	// Callback for Special Keys
	onTyped := func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			_, _ = process.Write([]byte{'\r'})
		}
	}

	// Callback for Char keys
	onRune := func(r rune) {
		_, _ = process.WriteString(string(r))
	}

	window.Canvas().SetOnTypedKey(onTyped)
	window.Canvas().SetOnTypedRune(onRune)

	buffer := [][]rune{}
	reader := bufio.NewReader(process)

	// Reading from pty
	go func() {
		line := []rune{}
		buffer = append(buffer, line)
		for {
			r, _, err := reader.ReadRune()

			if err != nil {
				if err == io.EOF {
					return
				}
				os.Exit(0)
			}

			line = append(line, r)
			buffer[len(buffer)-1] = line
			if r == '\n' {
				if len(buffer) > MaxBufferSize { // If the buffer is at capacity...
					buffer = buffer[1:] // ...pop the first line in the buffer
				}

				line = []rune{}
				buffer = append(buffer, line)
			}
		}
	}()

	// Updating the UI
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			textGrid.SetText("")
			var lines string
			for _, line := range buffer {
				lines = lines + string(line)
			}
			textGrid.SetText(string(lines))
		}
	}()

	// Display
	window.SetContent(
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(800, 600)),
			textGrid,
		),
	)

	window.ShowAndRun()
}
