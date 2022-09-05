package parser

import (
	"regexp"

	"github.com/hpcloud/tail"
	"github.com/juju/loggo"
)

var pkAuthRegexp = regexp.MustCompile("(chan_[0-9]+) -> PKAUTH (OPENPGP\\.3|[0-9A-F]+)$")
var pkSignRegexp = regexp.MustCompile("(chan_[0-9]+) <- PKSIGN$")
var eofRegexp = regexp.MustCompile("(chan_[0-9]+) <- \\[eof\\]$")
var okRegexp = regexp.MustCompile("(chan_[0-9]+) <- OK$")

type LogParser struct {
	logger         loggo.Logger
	parseEventChan chan<- *ParseEvent
}

type ParseEvent struct {
	ChanId string
	Type   string
}

func NewLogParser(logger loggo.Logger, parseEventChan chan<- *ParseEvent) *LogParser {
	return &LogParser{
		logger:         logger,
		parseEventChan: parseEventChan,
	}
}

func (p *LogParser) ParseLines(lines chan *tail.Line) {
	for line := range lines {
		p.parse(line)
	}
}

func (p *LogParser) parse(line *tail.Line) {
	p.logger.Child("log").Tracef(line.Text)

	if pkAuthMatch := pkAuthRegexp.FindSubmatch([]byte(line.Text)); pkAuthMatch != nil {
		p.logger.Infof("pkauth [%s]", pkAuthMatch[1])

		p.parseEventChan <- &ParseEvent{
			ChanId: string(pkAuthMatch[1]),
			Type:   "auth",
		}

	} else if pkSignMatch := pkSignRegexp.FindSubmatch([]byte(line.Text)); pkSignMatch != nil {
		p.logger.Infof("pksign [%s]", pkSignMatch[1])

		p.parseEventChan <- &ParseEvent{
			ChanId: string(pkSignMatch[1]),
			Type:   "sign",
		}

	} else if okMatch := okRegexp.FindSubmatch([]byte(line.Text)); okMatch != nil {
		p.logger.Debugf("ok [%s]", okMatch[1])

		p.parseEventChan <- &ParseEvent{
			ChanId: string(okMatch[1]),
			Type:   "ok",
		}
	} else if eofMatch := eofRegexp.FindSubmatch([]byte(line.Text)); eofMatch != nil {
		p.logger.Debugf("eof [%s]", eofMatch[1])

		p.parseEventChan <- &ParseEvent{
			ChanId: string(eofMatch[1]),
			Type:   "eof",
		}
	}
}
