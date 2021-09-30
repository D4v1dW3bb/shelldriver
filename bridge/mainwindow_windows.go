package bridge

//import "github.com/progrium/shelldriver/walk"
import "github.com/lxn/walk"

func init() {
	register(&MnWdw{})
}

type MnWdw struct {
	walk.MainWindow `mapstructure:",squash"`
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
