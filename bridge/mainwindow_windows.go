package bridge

import "github.com/Gipcomp/winapi"

func init() {
	register(&MnWdw{})
}

type MnWdw struct {
	winapi.MainWindow `mapstructure:",squash"`
}

func (m *MnWdw) Resource() interface{} {
	return &m.MainWindow
}

func (m *MnWdw) Discard() error {
	return nil
}

func (m *MnWdw) Apply() error {
	MainWindow = &m.MainWindow
	return nil
}
