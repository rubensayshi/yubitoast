package notifier

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gobuffalo/packr/v2"
	"github.com/juju/loggo"
)

type FyneNotifier struct {
	logger    loggo.Logger
	popupChan <-chan *PopupEvent
	box       *packr.Box
}

func NewFyneNotifier(logger loggo.Logger, popupChan <-chan *PopupEvent) Notifier {
	return &FyneNotifier{
		logger:    logger,
		popupChan: popupChan,
		box:       NewImgsBox(),
	}
}

func (n *FyneNotifier) Run() {
	fappLock := sync.RWMutex{}
	fappLock.Lock()

	var fapp fyne.App

	var focus string = ""

	iconBuf, err := n.box.Find("gnupg.png")
	if err != nil {
		panic(errors.Wrapf(err, "read gnupg.png"))
	}
	icon := fyne.NewStaticResource("gnupg.png", iconBuf)

	{
		var w fyne.Window
		pendingCnt := 0

		wait := make(chan bool, 1)

		go func() {
			for {
				select {
				case <-wait:
					<-time.After(500 * time.Millisecond)

				case event := <-n.popupChan:
					n.logger.Debugf("toggle: %v \n", event.Toggle)

					if !event.Toggle {
						pendingCnt--
					} else {
						pendingCnt++
					}

					// close previous open window
					if pendingCnt <= 0 && w != nil {
						n.logger.Debugf("force close")
						// will trigger OnClosed
						w.Close()
						// signal loop to wait a little bit
						wait <- true
						// reset counter to fix any odd bugs
						pendingCnt = 0
					}

					if pendingCnt > 0 {
						if w == nil {
							n.logger.Tracef("new! old focus: %s \n", focus)

							fappLock.RLock()
							// don't focus the window on show
							glfw.WindowHint(glfw.FocusOnShow, 0)
							// window is set to be always ontop
							glfw.WindowHint(glfw.Floating, 1)
							w = fapp.NewWindow("YubiToast")
							fappLock.RUnlock()

							w.Resize(fyne.NewSize(500, 0))
							w.SetIcon(icon)
							w.SetOnClosed(func() {
								n.logger.Child("gui").Debugf("on close")
								// @TODO: race
								w = nil
							})
						} else {
							n.logger.Tracef("reuse \n")
						}

						label := widget.NewLabel(fmt.Sprintf("Please touchy %d times!", pendingCnt))
						label.Alignment = fyne.TextAlignCenter

						btn := widget.NewButton("CLOSE", func() {
							n.logger.Child("gui").Debugf("click close")
							// will trigger OnClosed
							w.Close()
							// signal loop to wait a little bit
							wait <- true
						})
						btn.SetIcon(icon)

						w.SetContent(widget.NewVBox(
							label,
							btn,
						))

						// refocus window when event.Toggle=true
						if event.Toggle {
							w.Show()
						}
					}
				}
			}
		}()
	}

	// blocking, needs to be on main thread
	// stops blocking every time the last window is closed, so we do it forever
	for {
		n.logger.Debugf("RUN... \n")

		// make a new fapp
		fapp = app.New()

		// release the fappLock after 100ms so that fapp can be used
		//  in a go routine because fapp.Run() is blocking and needs to be on main thread...
		timer := time.AfterFunc(100*time.Millisecond, func() {
			fappLock.Unlock()
		})

		fapp.Run()

		n.logger.Debugf("RIP... \n")

		// quit to cleanup
		fapp.Quit()

		// cancel the unlock
		timer.Stop()

		// obtain lock to prevent further usage of fapp
		fappLock.Lock()

	}
}
