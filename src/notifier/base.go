package notifier

type Notifier interface {
	Run()
}

type PopupEvent struct {
	Toggle bool
}
