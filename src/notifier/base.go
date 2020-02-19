package notifier

import (
	"github.com/gobuffalo/packr/v2"
)

type Notifier interface {
	Run()
}

type PopupEvent struct {
	Toggle bool
}

func NewImgsBox() *packr.Box {
	box := packr.New("imgbox", "./imgs")

	return box
}
