package notifier

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/gobuffalo/packr"
	"github.com/juju/loggo"
	"github.com/rubensayshi/yubitoast/src/utils"
)

type FyneNotifier struct {
	logger    loggo.Logger
	popupChan <-chan *PopupEvent
	box       packr.Box
}

func NewFyneNotifier(logger loggo.Logger, popupChan <-chan *PopupEvent) Notifier {
	box := packr.NewBox(utils.ROOT + "/imgs")

	return &FyneNotifier{
		logger:    logger,
		popupChan: popupChan,
		box:       box,
	}
}

func (n *FyneNotifier) Run() {
	fapp := app.New()

	iconBuf, err := n.box.Find("gnupg.png")
	if err != nil {
		panic(errors.Wrapf(err, ""))
	}
	icon := fyne.NewStaticResource("gnupg.png", iconBuf)

	{
		var w fyne.Window
		pendingCnt := 0

		go func() {
			logger := n.logger.Child("view")
			for event := range n.popupChan {
				n.logger.Debugf("toggle: %v \n", event.Toggle)

				if !event.Toggle {
					pendingCnt--
				} else {
					pendingCnt++
				}

				// close previous open window
				if pendingCnt <= 0 && w != nil {
					n.logger.Debugf("force close")
					w.Close()
				}

				if pendingCnt > 0 {
					if w == nil {
						n.logger.Tracef("new \n")
						w = fapp.NewWindow("YubiToast")
						w.Resize(fyne.NewSize(500, 0))
						w.SetIcon(icon)
						w.SetOnClosed(func() {
							n.logger.Child("gui").Debugf("on close")
							// @TODO: race
							w = nil
							altTab(logger)
						})
					} else {
						n.logger.Tracef("reuse \n")
					}

					label := widget.NewLabel(fmt.Sprintf("Please touchy %d times!", pendingCnt))
					label.Alignment = fyne.TextAlignCenter

					btn := widget.NewButton("OK", func() {
						n.logger.Child("gui").Debugf("click close")
						w.Close()
						// @TODO: race
						w = nil
						altTab(logger)
					})
					btn.SetIcon(icon)

					w.SetContent(widget.NewVBox(
						label,
						btn,
					))

					// only (re)focus window when event.Toggle=true
					if event.Toggle {
						w.Show()
					}
				}
			}
		}()
	}

	// blocking, needs to be on main thread
	// stops blocking every time the last window is closed, so we do it forever
	for {
		n.logger.Debugf("Run... \n")
		fapp.Run()

		fapp = app.New()
	}
}

func altTab(logger loggo.Logger) {
	logger.Child("alttab").Debugf("alttab")
	//exec.Command("osascript", "-e", "tell application \"iTerm2\" to activate").Run()
	exec.Command("osascript", "-e", "tell application \"System Events\" to keystroke tab using command down").Run()
}
