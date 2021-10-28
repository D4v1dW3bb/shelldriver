package bridge

import (
	"unsafe"

	"github.com/Gipcomp/win32/user32"
	"github.com/Gipcomp/winapi"
	"github.com/progrium/shelldriver/shell"
)

// Indicator on mac is Systray on Winows
// Shelldrive originaly was designed for mac.
type Indicator struct {
	shell.Indicator `mapstructure:",squash"`

	// mainWindow *winapi.MainWindow
	target *winapi.NotifyIcon
	menu   *Menu
}

// Discard Discards the Indicator
func (i *Indicator) Discard() error {
	if i.menu != nil {
		i.menu.target.Dispose()
	}
	i.target.Dispose()
	return nil
}

// Apply applies the settings to the Indicator that is initiated
// in bridge.Main()
func (i *Indicator) Apply() error {
	// Retrieve our *NotifyIcon from the message window.
	ptr := user32.GetWindowLongPtr(MainWindow.Handle(), user32.GWLP_USERDATA)
	i.target = (*winapi.NotifyIcon)(unsafe.Pointer(ptr))

	// // Text
	// if i.Text != "" {
	// 	// Seams that text Icons do not work on windows
	// 	// Todo for a later moment: research if text Icons work
	// 	// or find a workaround for windows
	// 	icon, err := winapi.Resources.Icon(i.Text)
	// 	if err != nil {
	// 		return err
	// 	}

	// Icon
	if i.Icon != "" {

		icon, err := winapi.Resources.Icon(i.Icon)
		if err != nil {
			return err
		}

		err = i.target.SetIcon(icon)
		if err != nil {
			return err
		}
	}

	if i.Menu != nil {
		if i.menu == nil {
			i.menu = &Menu{}
		} else {
			i.target.ContextMenu().Dispose()
		}
		i.menu.ni = i.target
		i.menu.Menu = *i.Menu
		if err := i.menu.Apply(); err != nil {
			return err
		}
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := i.target.SetVisible(true); err != nil {
		return err
	}

	// Tooltip
	if i.menu.Tooltip != "" {
		if err := i.target.SetToolTip(i.menu.Tooltip); err != nil {
			return err
		}
	}

	return nil
}
