package notifier

import (
	"github.com/juju/loggo"
	"github.com/martinlindhe/notify"
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
	// create tmpfile holding the icon
	iconPath := createTmpIcon()
	n.logger.Tracef("icon: %s", iconPath)

	for event := range n.popupChan {
		n.logger.Debugf("toggle: %v \n", event.Toggle)

		if event.Toggle == true {
			notify.Notify(
				"Yubikey Touch Requested", "Yubikey Touch Requested",
				"Touchy YubiKey!", iconPath)
		}
	}
}
