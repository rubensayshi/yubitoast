package notifier

import (
	"io/ioutil"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"

	"github.com/juju/loggo"
	"github.com/martinlindhe/notify"
)

type NativeToastNotifier struct {
	logger    loggo.Logger
	popupChan <-chan *PopupEvent
	box       *packr.Box
}

func NewNativeToastNotifier(logger loggo.Logger, popupChan <-chan *PopupEvent) Notifier {
	return &NativeToastNotifier{
		logger:    logger,
		popupChan: popupChan,
		box:       NewImgsBox(),
	}
}
func makeTmpIcon(box *packr.Box) string {
	iconBuf, err := box.Find("gnupg.png")
	if err != nil {
		panic(errors.Wrapf(err, "read gnupg.png"))
	}

	tmpIcon, err := ioutil.TempFile("", "yubitoast-icon-gnupg")
	if err != nil {
		panic(errors.Wrapf(err, "make tmp icon"))
	}
	defer tmpIcon.Close()
	_, _ = tmpIcon.Write(iconBuf)

	return tmpIcon.Name()
}

func (n *NativeToastNotifier) Run() {
	tmpIcon := makeTmpIcon(n.box)
	n.logger.Debugf("icon: %s", tmpIcon)

	for event := range n.popupChan {
		n.logger.Debugf("toggle: %v \n", event.Toggle)

		if event.Toggle == true {
			notify.Notify("Yubikey Touch Requested", "", "Touchy YubiKey!", tmpIcon)
		}
	}
}
