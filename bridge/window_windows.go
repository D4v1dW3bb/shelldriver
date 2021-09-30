package bridge

import (
	"log"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/progrium/shelldriver/shell"
	//"github.com/progrium/shelldriver/walk"
)

func init() {
	register(&Window{})
}

type Window struct {
	shell.Window `mapstructure:",squash"`

	mainWindow  *walk.MainWindow
	windowStyle int
	// // windowExstyle prepaired for WS_EX_..... style messaging
	// // see also functions at the end of this file
	// windowExStyle int32
}

func (w *Window) Resource() interface{} {
	return &w.Window
}

func (w *Window) Discard() error {
	// TODO: image, webview
	win.SendMessage(w.mainWindow.Handle(), win.WM_CLOSE, 0, 0)
	win.DestroyWindow(w.mainWindow.Handle())
	return nil
}

func (w *Window) Apply() error {
	if MainWindow == nil {
		for MainWindow == nil {
			time.Sleep(10 * time.Millisecond)
		}
	}

	w.mainWindow = MainWindow

	w.mainWindow.SetTitle(w.Title)

	// Center
	if w.Center {
		var rect win.RECT
		screen := win.GetDesktopWindow()
		win.GetClientRect(screen, &rect)

		//fmt.Println(win.GetDeviceCaps(hdc, win.HORZSIZE))
		w.Position.X = float64(rect.Right/2) - (w.Size.W / 2)
		w.Position.Y = float64(rect.Bottom/2) - (w.Size.H / 2)
	}

	// AlwaysOnTop
	if w.AlwaysOnTop {
		win.SetWindowPos(w.mainWindow.Handle(), win.HWND_TOPMOST, int32(w.Position.X), int32(w.Position.Y), int32(w.Size.W), int32(w.Size.H), win.SWP_SHOWWINDOW)
	} else {
		win.SetWindowPos(w.mainWindow.Handle(), win.HWND_NOTOPMOST, int32(w.Position.X), int32(w.Position.Y), int32(w.Size.W), int32(w.Size.H), win.SWP_SHOWWINDOW)
	}

	// URL
	if w.URL != "" {
		WebView.SetURL(w.URL)
	} else {
		WebView.SetURL("")
	}

	if w.Closable {
		closable(win.MF_BYCOMMAND | win.MF_ENABLED)
	} else {
		closable(win.MF_BYCOMMAND | win.MF_DISABLED | win.MF_GRAYED)
	}

	// Get current window style
	w.windowStyle = int(win.GetWindowLong(w.mainWindow.Handle(), win.GWL_STYLE))

	// Minimizable
	if w.Minimizable && !w.checkWindowStyle(win.WS_MINIMIZEBOX) {
		w.addWindowStyle(win.WS_MINIMIZEBOX)
	} else if !w.Minimizable && w.checkWindowStyle(win.WS_MINIMIZEBOX) {
		w.removeWindowStyle(win.WS_MINIMIZEBOX)
	}

	// Resizable
	if w.Resizable && !w.checkWindowStyle(win.WS_SIZEBOX) {
		w.addWindowStyle(win.WS_SIZEBOX | win.WS_MAXIMIZEBOX)
	} else if !w.Resizable && w.checkWindowStyle(win.WS_SIZEBOX) {
		w.removeWindowStyle(win.WS_SIZEBOX | win.WS_MAXIMIZEBOX)
	}

	// Borderless
	needsTitleBar := w.Closable || w.Minimizable
	if w.Borderless {
		if needsTitleBar {
			w.removeWindowStyle(win.WS_BORDER)
			w.addWindowStyle(win.WS_MAXIMIZEBOX)
		} else {
			w.removeWindowStyle(win.WS_BORDER | win.WS_DLGFRAME | win.WS_THICKFRAME | win.WS_SYSMENU | win.WS_MAXIMIZEBOX)
		}

	} else {
		if needsTitleBar {
			w.addWindowStyle(win.WS_BORDER | win.WS_DLGFRAME | win.WS_THICKFRAME | win.WS_SYSMENU | win.WS_MAXIMIZEBOX)
		} else {
			w.addWindowStyle(win.WS_BORDER)
			w.removeWindowStyle(win.WS_DLGFRAME | win.WS_THICKFRAME | win.WS_SYSMENU | win.WS_MAXIMIZEBOX)
		}
	}
	setWindowStyle(int32(w.windowStyle))

	// Background
	// Todo: background of WebView
	if w.Background != nil {
		brush, _ := walk.NewSolidColorBrush(walk.RGB(byte(w.Background.R), byte(w.Background.G), byte(w.Background.B)))
		w.mainWindow.SetBackground(brush)
	}

	// Image
	// Todo: image placement and size, stretch, center...
	if w.Image != "" {
		bmp, err := walk.NewBitmapFromFileForDPI(w.Image, w.mainWindow.DPI())
		if err != nil {
			log.Panicln(err)
		}
		bmb, err := walk.NewBitmapBrush(bmp)
		if err != nil {
			log.Panicln(err)
		}
		w.mainWindow.SetBackground(bmb)
	}

	// Todo IgnoreMouse needs hook in walk to accomplish this
	// if w.IgnoreMouse {

	// } else {

	// }

	return nil

	// Todo: Movable without Title
	// if w.Title != "" {
	// 	w.target.SetTitle(w.Title)
	// } else {
	// 	w.target.SetMovableByWindowBackground(true)
	// 	w.target.SetTitlebarAppearsTransparent(true)
	// }

	// Todo: CornerRadius
	// 	v.Layer().SetCornerRadius(w.CornerRadius)

	// Todo: IgnoreMouse
	// if w.IgnoreMouse {
	// 	w.target.SetIgnoresMouseEvents(true)
	// }

}

func (w *Window) addWindowStyle(style int) {
	w.windowStyle = w.windowStyle | style
}

func (w *Window) removeWindowStyle(style int) {
	w.windowStyle = w.windowStyle &^ style
}

func (w *Window) checkWindowStyle(dwType int32) bool {
	return ((win.GetWindowLong(w.mainWindow.Handle(), win.GWL_STYLE) & dwType) != 0)
}

func closable(style uint32) {
	win.EnableMenuItem(
		win.GetSystemMenu(
			MainWindow.Handle(),
			false),
		win.SC_CLOSE,
		style)
}

func setWindowStyle(style int32) {
	win.SetWindowLong(
		MainWindow.Handle(),
		win.GWL_STYLE,
		style,
	)
}

// // WS_EX... style messages functions for windows
// func setWindowExStyle(style int32) {
// 	win.SetWindowLong(
// 		MainWindow.Handle(),
// 		win.GWL_EXSTYLE,
// 		style,
// 	)
// }

// func (w *Window) addWindowExStyle(style int32) {
// 	w.windowExStyle = w.windowExStyle | style
// }

// func (w *Window) removeWindowExStyle(style int32) {
// 	w.windowExStyle = w.windowExStyle &^ style
// }
