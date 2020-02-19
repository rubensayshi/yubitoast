package toaster

import (
	"time"

	"github.com/hpcloud/tail"
	"github.com/juju/loggo"
	"github.com/rubensayshi/yubitoast/src/notifier"
	"github.com/rubensayshi/yubitoast/src/parser"
)

type Toaster struct {
	logger   loggo.Logger
	notifier notifier.Notifier
	parser   *parser.LogParser

	parseEventChan <-chan *parser.ParseEvent
	popupChan      chan<- *notifier.PopupEvent

	pending  map[string]string
	timeouts map[string]*time.Timer
}

func NewToaster(
	logger loggo.Logger,
	notifierFn func(loggo.Logger, <-chan *notifier.PopupEvent) notifier.Notifier) *Toaster {

	popupChan := make(chan *notifier.PopupEvent)
	parseEventChan := make(chan *parser.ParseEvent)

	return &Toaster{
		logger:   logger.Child("toaster"),
		notifier: notifierFn(logger.Child("notifier"), popupChan),
		parser:   parser.NewLogParser(logger.Child("parser"), parseEventChan),

		parseEventChan: parseEventChan,
		popupChan:      popupChan,

		pending:  make(map[string]string),
		timeouts: make(map[string]*time.Timer),
	}
}

func (t *Toaster) Run(lines chan *tail.Line) {
	go t.parser.ParseLines(lines)

	go t.run()

	// blocking, needs to be on main thread!
	t.notifier.Run()
}

func (t *Toaster) run() {
	for parseEvent := range t.parseEventChan {
		switch parseEvent.Type {
		case "auth":
			fallthrough
		case "sign":
			t.pending[parseEvent.ChanId] = parseEvent.Type
			t.timeouts[parseEvent.ChanId] = time.AfterFunc(15*time.Second, func() {
				delete(t.timeouts, parseEvent.ChanId)
				delete(t.pending, parseEvent.ChanId)

				t.popupChan <- &notifier.PopupEvent{Toggle: false}
			})

			t.popupChan <- &notifier.PopupEvent{Toggle: true}

		case "ok":
			if t.pending[parseEvent.ChanId] == "auth" {
				t.timeouts[parseEvent.ChanId].Stop()
				delete(t.timeouts, parseEvent.ChanId)
				delete(t.pending, parseEvent.ChanId)

				t.popupChan <- &notifier.PopupEvent{Toggle: false}
			}
		case "eof":
			if t.pending[parseEvent.ChanId] == "sign" {
				t.timeouts[parseEvent.ChanId].Stop()
				delete(t.timeouts, parseEvent.ChanId)
				delete(t.pending, parseEvent.ChanId)

				t.popupChan <- &notifier.PopupEvent{Toggle: false}
			}

		}
	}
}
