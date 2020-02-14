package notifier

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gobuffalo/packr/v2"
	"github.com/juju/loggo"
	"github.com/pkg/errors"
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

					if pendingCnt > 0 {
						// first popup we need to create a new window
						if w == nil {
							n.logger.Debugf("make a new window")

							// we need to get a read lock on the fapp
							fappLock.RLock()
							// don't focus the window on show
							glfw.WindowHint(glfw.FocusOnShow, 0)
							// window is set to be always ontop
							glfw.WindowHint(glfw.Floating, 1)
							w = fapp.NewWindow("YubiToast")

							// set a CloseIntercept to Hide instead of Close
							//  because Close results in the app destructing
							w.SetCloseIntercept(func() {
								w.Hide()
							})
							w.Resize(fyne.NewSize(500, 100))
							w.SetIcon(icon)
							w.SetOnClosed(func() {
								n.logger.Child("gui").Debugf("OnClosed")
							})

							fappLock.RUnlock()

						} else {
							n.logger.Debugf("reuse window")
						}

						n.logger.Debugf("add elements to window ...")

						label := widget.NewLabel(fmt.Sprintf("Please touchy %d times!", pendingCnt))
						label.Alignment = fyne.TextAlignCenter

						btn := widget.NewButton("CLOSE", func() {
							n.logger.Child("gui").Debugf("close button clicked")
							// will trigger OnClosed
							w.Hide()
							// signal loop to wait a little bit
							wait <- true
						})
						btn.SetIcon(icon)

						w.SetContent(container.NewVBox(
							label,
							btn,
						))

						n.logger.Debugf("show")
						// refocus window when event.Toggle=true
						if event.Toggle {
							w.Show()
						}
					}
				}
			}
		}()
	}

	// make a new fapp
	fapp = app.New()

	// release the fappLock after 100ms so that fapp can be used
	//  in a go routine because fapp.Run() is blocking and needs to be on main thread...
	time.AfterFunc(100*time.Millisecond, func() {
		fappLock.Unlock()
	})

	// blocking, needs to be on main thread
	fapp.Run()

	// quit to cleanup
	fapp.Quit()

}
