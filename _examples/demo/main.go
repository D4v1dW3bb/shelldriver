package main

import (
	"log"
	"os"

	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/shelldriver/shell"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sh := shell.New(nil)
	sh.Debug = os.Stderr
	must(sh.Open())
	defer sh.Close()
	w := shell.Window{
		Title:    "Demo",
		Size:     shell.Size{W: 800, H: 400},
		Position: shell.Point{X: 200, Y: 200},
		// Center:      true,
		// Borderless: true,
		Minimizable: true,
		Closable:    true,
		Resizable:   true,
		// AlwaysOnTop: true,
		// URL: "www.google.com",
		// Webview background change has to be done yet
		// Background: &shell.Color{
		// 	R: 128,
		// 	B: 128,
		// 	G: 255,
		// },
		Image: "./resources/tractor.jpg",
	}
	must(sh.Sync(&w))

	i := shell.Indicator{
		// Text Icons are not supported yet on windows
		// You can only use it on mac.
		// The little tractor is possible a font icon
		// Text: "🚜",
		Icon: "./resources/tractor.ico",
		Menu: &shell.Menu{
			Items: []shell.MenuItem{
				{Title: "Window", Enabled: true, SubItems: []shell.MenuItem{
					{Title: "Step", Enabled: true, OnClick: fn.Callback(func() {
						w.Position.X += 20
						w.Position.Y += 20
						must(sh.Sync(&w))
					})},
					{Separator: true},
					{Title: "Always On Top",
						Enabled: true,
						// Checkable added for windows functionality
						Checkable: true,
						// TODO Mac: Always On Top check/state
						Checked: w.AlwaysOnTop,
						OnClick: fn.Callback(func() {
							w.AlwaysOnTop = !w.AlwaysOnTop
							log.Println("On top? ", w.AlwaysOnTop)
							must(sh.Sync(&w))
						})},
				}},
				{Separator: true},
				{Title: "Quit", Enabled: true},
			},
		},
	}

	must(sh.Sync(&i))

	if err := sh.Wait(); err != nil {
		log.Fatal(err)
	}
}
