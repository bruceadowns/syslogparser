package rfc3164

import (
	"bytes"
	"strconv"
	"time"

	"github.com/bruceadowns/syslogparser"
)

// Parser structure
type Parser struct {
	buff     []byte
	cursor   int
	l        int
	priority syslogparser.Priority
	version  int
	header   header
	message  rfc3164message
	hostname string
}

type header struct {
	timestamp time.Time
	hostname  string
}

type rfc3164message struct {
	app     string
	pid     string
	content string
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		buff:   buff,
		cursor: 0,
		l:      len(buff),
	}
}

// Hostname ...
func (p *Parser) Hostname(hostname string) {
	p.hostname = hostname
}

// Parse ...
func (p *Parser) Parse() error {
	pri, err := p.parsePriority()
	if err != nil {
		return err
	}

	hdr, err := p.parseHeader()
	if err != nil {
		return err
	}

	if p.buff[p.cursor] == ' ' {
		p.cursor++
	}

	msg, err := p.parseMessage()
	if err != syslogparser.ErrEOL {
		return err
	}

	p.priority = pri
	p.version = syslogparser.NoVersion
	p.header = hdr
	p.message = msg

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	return syslogparser.LogParts{
		"timestamp": syslogparser.Epoch(p.header.timestamp),
		"hostname":  p.header.hostname,
		"app_name":  p.message.app,
		"proc_id":   p.message.pid,
		"content":   p.message.content,
		"priority":  strconv.Itoa(p.priority.P),
		"facility":  strconv.Itoa(p.priority.F.Value),
		"severity":  strconv.Itoa(p.priority.S.Value),
	}
}

func (p *Parser) parsePriority() (syslogparser.Priority, error) {
	return syslogparser.ParsePriority(p.buff, &p.cursor, p.l)
}

func (p *Parser) parseHeader() (header, error) {
	hdr := header{}
	var err error

	ts, err := p.parseTimestamp()
	if err != nil {
		return hdr, err
	}

	hostname, err := p.parseHostname()
	if err != nil {
		return hdr, err
	}

	hdr.timestamp = ts
	hdr.hostname = hostname

	return hdr, nil
}

func (p *Parser) parseMessage() (rfc3164message, error) {
	msg := rfc3164message{}
	var err error

	app, pid, err := p.parseTag()
	if err != nil {
		return msg, err
	}

	msg.app = app
	msg.pid = pid

	content, err := p.parseContent()
	if err != syslogparser.ErrEOL {
		return msg, err
	}

	msg.content = content

	return msg, err
}

// https://tools.ietf.org/html/rfc3164#section-4.1.2
func (p *Parser) parseTimestamp() (time.Time, error) {
	var ts time.Time
	var err error
	var tsFmtLen int
	var sub []byte

	tsFmts := []string{
		"Jan 02 15:04:05",
		"Jan  2 15:04:05",
	}

	found := false
	for _, tsFmt := range tsFmts {
		tsFmtLen = len(tsFmt)

		if p.cursor+tsFmtLen > p.l {
			continue
		}

		sub = p.buff[p.cursor : tsFmtLen+p.cursor]
		ts, err = time.ParseInLocation(tsFmt, string(sub), time.UTC)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		p.cursor = tsFmtLen

		// XXX : If the timestamp is invalid we try to push the cursor one byte
		// XXX : further, in case it is a space
		if p.cursor < p.l && p.buff[p.cursor] == ' ' {
			p.cursor++
		}

		return ts, syslogparser.ErrTimestampUnknownFormat
	}

	fixTimestampIfNeeded(&ts)

	p.cursor += tsFmtLen

	if p.cursor < p.l && p.buff[p.cursor] == ' ' {
		p.cursor++
	}

	return ts, nil
}

func (p *Parser) parseHostname() (string, error) {
	if p.hostname != "" {
		return p.hostname, nil
	}

	return syslogparser.ParseHostname(p.buff, &p.cursor, p.l)
}

// http://tools.ietf.org/html/rfc3164#section-4.1.3
func (p *Parser) parseTag() (string, string, error) {
	var b byte
	var endOfTag, endOfApp, endOfPid bool
	var app, pid []byte
	var foundApp, foundPid bool

	from := p.cursor

	for {
		b = p.buff[p.cursor]

		endOfApp = b == '['
		endOfPid = b == ']'
		endOfTag = (b == ':' || b == ' ')

		if endOfApp {
			app = p.buff[from:p.cursor]
			from = p.cursor
			foundApp = true
		} else if endOfPid {
			if !foundApp {
				app = p.buff[from:p.cursor]
				foundApp = true
			} else if !foundPid {
				pid = p.buff[from+1 : p.cursor]
				foundPid = true
			}

			p.cursor++
			p.cursor++
			break
		} else if endOfTag {
			if !foundApp {
				app = p.buff[from:p.cursor]
				foundApp = true
			} else if !foundPid {
				pid = p.buff[from+1 : p.cursor]
				foundPid = true
			}

			p.cursor++
			break
		}

		p.cursor++
	}

	if p.cursor < p.l && p.buff[p.cursor] == ' ' {
		p.cursor++
	}

	return string(app), string(pid), nil
}

func (p *Parser) parseContent() (string, error) {
	if p.cursor > p.l {
		return "", syslogparser.ErrEOL
	}

	content := bytes.Trim(p.buff[p.cursor:p.l], " ")
	p.cursor += len(content)

	return string(content), syslogparser.ErrEOL
}

func fixTimestampIfNeeded(ts *time.Time) {
	y := ts.Year()
	if y == 0 {
		y = time.Now().Year()
	}

	newTs := time.Date(y, ts.Month(), ts.Day(), ts.Hour(), ts.Minute(),
		ts.Second(), ts.Nanosecond(), ts.Location())

	*ts = newTs
}
