package notifier

import (
	"github.com/juju/loggo"
	"github.com/martinlindhe/notify"
	"github.com/rubensayshi/yubitoast/src/utils"
)

type NativeToastNotifier struct {
	logger    loggo.Logger
	popupChan <-chan *PopupEvent
}

func NewNativeToastNotifier(logger loggo.Logger, popupChan <-chan *PopupEvent) Notifier {
	return &NativeToastNotifier{
		logger:    logger,
		popupChan: popupChan,
	}
}

func (n *NativeToastNotifier) Run() {
	for event := range n.popupChan {
		n.logger.Debugf("toggle: %v \n", event.Toggle)

		if event.Toggle == true {
			notify.Notify("Yubikey Touch Requested", "", "Touchy YubiKey!", utils.ROOT+"/imgs/gnupg.png")
		}
	}
}
