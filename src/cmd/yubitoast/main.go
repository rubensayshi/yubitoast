package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	"github.com/hpcloud/tail"
	"github.com/martinlindhe/notify"
)

var pkAuthRegexp = regexp.MustCompile("PKAUTH OPENPGP\\.3$")
var pkSignRegexp = regexp.MustCompile("PKSIGN --hash=.+ OPENPGP\\.1$")

var logfile = flag.String("logfile", "/var/log/gpg-agent.log", "path to gpg-agent.log")

func main() {
	flag.Parse()

	t, err := tail.TailFile(*logfile, tail.Config{
		Follow: true,
		// seek to end of file
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
	})
	if err != nil {
		panic(errors.Wrapf(err, "tail gpg-agent log [%s]", (logfile)))
	}

	// now we wait
	for line := range t.Lines {
		fmt.Println(line.Text)

		if pkAuthRegexp.Match([]byte(line.Text)) {
			showNotification("Authenticate")
		} else if pkSignRegexp.Match([]byte(line.Text)) {
			showNotification("Sign")
		}
	}
}

func showNotification(msg string) {
	notify.Notify("Yubikey Touch Requested", "", msg, "gnupg.png")
}
