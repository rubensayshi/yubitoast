package main

import (
	"flag"

	"github.com/hpcloud/tail"
	"github.com/juju/loggo"
	"github.com/pkg/errors"
	"github.com/rubensayshi/yubitoast/src/log"
	"github.com/rubensayshi/yubitoast/src/notifier"
	"github.com/rubensayshi/yubitoast/src/toaster"
)

var fLogfile = flag.String("logfile", "/var/log/gpg-agent.log", "path to gpg-agent.log")
var fDebug = flag.Bool("debug", false, "verbose / debug logging")
var fTrace = flag.Bool("trace", false, "super verbose / trace logging")
var fNotifier = flag.String("notifier", "fyne", "notifier type to use; [native | fyne]")

func main() {
	flag.Parse()

	logger := loggo.GetLogger("yubitoast")
	logger.SetLogLevel(loggo.INFO)
	if *fTrace {
		logger.SetLogLevel(loggo.TRACE)
	} else if *fDebug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	t, err := tail.TailFile(*fLogfile, tail.Config{
		Follow: true,
		// seek to end of file
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		Logger: log.TailLoggoLogger{logger.Child("tail")},
	})
	if err != nil {
		panic(errors.Wrapf(err, "tail gpg-agent log [%s]", (fLogfile)))
	}

	var notifierFn func(loggo.Logger, <-chan *notifier.PopupEvent) notifier.Notifier
	switch *fNotifier {
	case "native":
		notifierFn = notifier.NewNativeToastNotifier
	case "fyne":
		notifierFn = notifier.NewFyneNotifier
	default:
		panic(errors.Errorf("Invalid notifier: %s", *fNotifier))
	}

	toaster := toaster.NewToaster(logger, notifierFn)
	toaster.Run(t.Lines)
}
