package bridge

import (
	"log"
	"strings"

	"github.com/Gipcomp/winapi"
	"github.com/Gipcomp/winapi/declarative"
)

// MainWindow Holds the main window
var MainWindow *winapi.MainWindow

// WebView holds the webview child
var WebView *winapi.WebView

// IN holds the notification icon
var IN *winapi.NotifyIcon

// Main initializes the mainwindow, webview and notification icon
func Main() {
	declarative.MainWindow{
		AssignTo: &MainWindow,
		Title:    "ShellDriverMain",
		MinSize:  declarative.Size{Height: 100, Width: 100},
		Size:     declarative.Size{Height: 100, Width: 300},
		Visible:  false,
		Layout:   declarative.VBox{},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}
				return "stop", nil
			},
		},
	}.Create()

	// Enable this to enable webview on windows and to use w.URL.
	WebView, _ = winapi.NewWebView(MainWindow)
	WebView.SetName("Webview")

	IN, err := winapi.NewNotifyIcon(MainWindow)
	if err != nil {
		log.Panic(err)
	}
	MainWindow.Run()
	IN.Dispose()
}
