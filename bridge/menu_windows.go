package bridge

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lxn/walk"
	"github.com/progrium/shelldriver/shell"
	//"github.com/progrium/shelldriver/walk"
)

type Menu struct {
	shell.Menu `mapstructure:",squash"`

	ni     *walk.NotifyIcon
	target *walk.Menu
}

func (m *Menu) Resource() interface{} {
	return &m.Menu
}

func (m *Menu) Discard() error {
	m.target.Dispose()
	return nil
}

func (m *Menu) Apply() error {
	if m.target != nil {
		m.target.Dispose()
	}

	for _, i := range m.Items {
		if i.SubItems == nil {
			m.ni.ContextMenu().Actions().Add(MenuItem(i))
		} else {
			subMenu, _ := walk.NewMenu()
			for _, si := range i.SubItems {
				subMenu.Actions().Add(MenuItem(si))
			}
			action, _ := m.ni.ContextMenu().Actions().AddMenu(subMenu)
			action.SetText(i.Title)
			action.SetEnabled(i.Enabled)
			action.SetToolTip(i.Tooltip)
			// Checked
			if i.Checked {
				action.SetChecked(i.Checked)

			}

			// Icon
			// todo: test if this has the same effect as setting an Icon on Mac

			if i.Icon != "" {
				icon, err := walk.Resources.Image(i.Icon)
				if err == nil {
					action.SetImage(icon)
				}
			}
		}

	}
	return nil
}

func MenuItem(i shell.MenuItem) *walk.Action {
	obj := walk.NewAction()
	// Separator
	if i.Separator {
		action := walk.NewSeparatorAction()
		return action
	}

	// Title
	obj.SetText(i.Title)

	// Enabled
	obj.SetEnabled(i.Enabled)

	// Tooltip
	obj.SetToolTip(i.Tooltip)

	// Checked
	if i.Checked {

		obj.SetChecked(i.Checked)

	}

	// Icon
	// todo: test if this has the same effect as setting an Icon on Mac

	if i.Icon != "" {
		icon, err := walk.Resources.Image(i.Icon)
		if err == nil {
			obj.SetImage(icon)
		}
	}

	// Quit default action
	if i.Title == "Quit" {
		obj.SetText("Quit")
		obj.Triggered().Attach(func() { walk.App().Exit(0) })
	}

	if i.Checkable {
		obj.SetCheckable(true)
	}

	// Todo: OnClick action check/unchack
	// Figureout how to handle the OnClick action on Windows
	if i.OnClick != nil {
		obj.Triggered().Attach(func() {
			go func() {
				if obj.Checkable() {
					obj.SetChecked(obj.Checked())
				}
				ff := *i.OnClick
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				_, err := ff.Call(ctx, nil, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "remote callback: %v\n", err)
				}
			}()
		})
	}

	// SubItems
	if len(i.SubItems) > 0 {
		var sub *walk.Menu

		for _, i := range i.SubItems {
			sub.Actions().Add(MenuItem(i))
			log.Println("Added: ", i)
		}
		obj = walk.NewMenuAction(sub)
		if i.Title != "" {
			obj.SetText(i.Title)
		}

		if i.Icon != "" {
			icon, err := walk.Resources.Image(i.Icon)
			if err == nil {
				obj.SetImage(icon)
			}
		}
	}

	return obj
}
