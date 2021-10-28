package bridge

import (
	"log"
	"time"

	"github.com/Gipcomp/win32/gdi32"
	"github.com/Gipcomp/win32/user32"
	"github.com/Gipcomp/win32/winuser"
	"github.com/Gipcomp/winapi"
	"github.com/progrium/shelldriver/shell"
)

func init() {
	register(&Window{})
}

type Window struct {
	shell.Window `mapstructure:",squash"`

	mainWindow  *winapi.MainWindow
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
	user32.SendMessage(w.mainWindow.Handle(), user32.WM_CLOSE, 0, 0)
	user32.DestroyWindow(w.mainWindow.Handle())
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
		var rect gdi32.RECT
		screen := user32.GetDesktopWindow()
		user32.GetClientRect(screen, &rect)

		//fmt.Println(user32.GetDeviceCaps(hdc, user32.HORZSIZE))
		w.Position.X = float64(rect.Right/2) - (w.Size.W / 2)
		w.Position.Y = float64(rect.Bottom/2) - (w.Size.H / 2)
	}

	// AlwaysOnTop
	if w.AlwaysOnTop {
		user32.SetWindowPos(w.mainWindow.Handle(), user32.HWND_TOPMOST, int32(w.Position.X), int32(w.Position.Y), int32(w.Size.W), int32(w.Size.H), user32.SWP_SHOWWINDOW)
	} else {
		user32.SetWindowPos(w.mainWindow.Handle(), user32.HWND_NOTOPMOST, int32(w.Position.X), int32(w.Position.Y), int32(w.Size.W), int32(w.Size.H), user32.SWP_SHOWWINDOW)
	}

	// URL
	if w.URL != "" {
		WebView.SetURL(w.URL)
	} else {
		WebView.SetURL("")
	}

	if w.Closable {
		closable(winuser.MF_BYCOMMAND | winuser.MF_ENABLED)
	} else {
		closable(winuser.MF_BYCOMMAND | winuser.MF_DISABLED | winuser.MF_GRAYED)
	}

	// Get current window style
	w.windowStyle = int(user32.GetWindowLong(w.mainWindow.Handle(), user32.GWL_STYLE))

	// Minimizable
	if w.Minimizable && !w.checkWindowStyle(user32.WS_MINIMIZEBOX) {
		w.addWindowStyle(user32.WS_MINIMIZEBOX)
	} else if !w.Minimizable && w.checkWindowStyle(user32.WS_MINIMIZEBOX) {
		w.removeWindowStyle(user32.WS_MINIMIZEBOX)
	}

	// Resizable
	if w.Resizable && !w.checkWindowStyle(user32.WS_SIZEBOX) {
		w.addWindowStyle(user32.WS_SIZEBOX | user32.WS_MAXIMIZEBOX)
	} else if !w.Resizable && w.checkWindowStyle(user32.WS_SIZEBOX) {
		w.removeWindowStyle(user32.WS_SIZEBOX | user32.WS_MAXIMIZEBOX)
	}

	// Borderless
	needsTitleBar := w.Closable || w.Minimizable
	if w.Borderless {
		if needsTitleBar {
			w.removeWindowStyle(user32.WS_BORDER)
			w.addWindowStyle(user32.WS_MAXIMIZEBOX)
		} else {
			w.removeWindowStyle(user32.WS_BORDER | user32.WS_DLGFRAME | user32.WS_THICKFRAME | user32.WS_SYSMENU | user32.WS_MAXIMIZEBOX)
		}

	} else {
		if needsTitleBar {
			w.addWindowStyle(user32.WS_BORDER | user32.WS_DLGFRAME | user32.WS_THICKFRAME | user32.WS_SYSMENU | user32.WS_MAXIMIZEBOX)
		} else {
			w.addWindowStyle(user32.WS_BORDER)
			w.removeWindowStyle(user32.WS_DLGFRAME | user32.WS_THICKFRAME | user32.WS_SYSMENU | user32.WS_MAXIMIZEBOX)
		}
	}
	setWindowStyle(int32(w.windowStyle))

	// Background
	// Todo: background of WebView
	if w.Background != nil {
		brush, _ := winapi.NewSolidColorBrush(winapi.RGB(byte(w.Background.R), byte(w.Background.G), byte(w.Background.B)))
		w.mainWindow.SetBackground(brush)
	}

	// Image
	// Todo: image placement and size, stretch, center...
	if w.Image != "" {
		bmp, err := winapi.NewBitmapFromFileForDPI(w.Image, w.mainWindow.DPI())
		if err != nil {
			log.Panicln(err)
		}
		bmb, err := winapi.NewBitmapBrush(bmp)
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
	return ((user32.GetWindowLong(w.mainWindow.Handle(), user32.GWL_STYLE) & dwType) != 0)
}

func closable(style uint32) {
	user32.EnableMenuItem(
		user32.GetSystemMenu(
			MainWindow.Handle(),
			false),
		user32.SC_CLOSE,
		style)
}

func setWindowStyle(style int32) {
	user32.SetWindowLong(
		MainWindow.Handle(),
		user32.GWL_STYLE,
		style,
	)
}

// // WS_EX... style messages functions for windows
// func setWindowExStyle(style int32) {
// 	user32.SetWindowLong(
// 		MainWindow.Handle(),
// 		user32.GWL_EXSTYLE,
// 		style,
// 	)
// }

// func (w *Window) addWindowExStyle(style int32) {
// 	w.windowExStyle = w.windowExStyle | style
// }

// func (w *Window) removeWindowExStyle(style int32) {
// 	w.windowExStyle = w.windowExStyle &^ style
// }
