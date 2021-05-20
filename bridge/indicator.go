package bridge

import (
	"encoding/base64"

	"github.com/progrium/macbridge/handle"
	"github.com/progrium/macbridge/shell"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type Indicator struct {
	shell.Indicator
}

func (i *Indicator) Resource() (*handle.Handle, interface{}) {
	return &i.Handle, &i.Indicator
}

func (s *Indicator) Apply(target objc.Object) (objc.Object, error) {
	obj := cocoa.NSStatusItem{Object: target}
	if target == nil {
		obj = cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		target = obj.Object
	}
	obj.Button().SetTitle(s.Text)
	if s.Icon != "" {
		b, err := base64.StdEncoding.DecodeString(s.Icon)
		if err != nil {
			return nil, err
		}
		data := core.NSData_WithBytes(b, uint64(len(b)))
		image := cocoa.NSImage_InitWithData(data)
		image.SetSize(core.Size(16.0, 16.0))
		image.SetTemplate(true)
		obj.Button().SetImage(image)
		if s.Text != "" {
			obj.Button().SetImagePosition(cocoa.NSImageLeft)
		} else {
			obj.Button().SetImagePosition(cocoa.NSImageOnly)
		}

	}
	if s.Menu != nil {
		var menu objc.Object
		var err error
		m := &Menu{*s.Menu}
		if menu, err = m.Apply(menu); err != nil {
			return nil, err
		}
		obj.SetMenu(cocoa.NSMenu{Object: menu})
	}
	return target, nil
}
