package bridge

func init() {
	register(&Indicator{})
}

func (i *Indicator) Resource() interface{} {
	return &i.Indicator
}
