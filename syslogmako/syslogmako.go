package syslogmako

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/bruceadowns/syslogparser"
	"github.com/bruceadowns/syslogparser/mako"
)

// Parser structure
type Parser struct {
	buff     []byte
	cursor   int
	l        int
	version  int
	header   header
	message  rfc3164Message
	location *time.Location
	hostname string
}

type header struct {
	timestamp time.Time
	hostname  string
}

type rfc3164Message struct {
	app     string
	pid     string
	content string
	mako    mako.JSON
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		buff:     buff,
		cursor:   0,
		l:        len(buff),
		location: time.UTC,
	}
}

// Location ...
func (p *Parser) Location(location *time.Location) {
	p.location = location
}

// Hostname ...
func (p *Parser) Hostname(hostname string) {
	p.hostname = hostname
}

// Parse ...
func (p *Parser) Parse() error {
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

	p.version = syslogparser.NoVersion
	p.header = hdr
	p.message = msg

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	return syslogparser.LogParts{
		"timestamp": p.header.timestamp,
		"hostname":  p.header.hostname,
		"app":       p.message.app,
		"pid":       p.message.pid,
		"content":   p.message.content,
		"mako":      p.message.mako,
	}
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

func (p *Parser) parseMessage() (rfc3164Message, error) {
	msg := rfc3164Message{}
	var err error

	app, pid, err := p.parseApp()
	if err != nil {
		return msg, err
	}

	content, err := p.parseContent()
	if err != syslogparser.ErrEOL {
		return msg, err
	}

	msg.app = app
	msg.pid = pid
	msg.content = content

	if err := json.NewDecoder(bytes.NewBufferString(msg.content)).Decode(&msg.mako); err != nil {
		return msg, err
	}

	return msg, syslogparser.ErrEOL
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
		ts, err = time.ParseInLocation(tsFmt, string(sub), p.location)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		p.cursor = tsFmtLen

		// XXX : If the timestamp is invalid we try to push the cursor one byte
		// XXX : further, in case it is a space
		if (p.cursor < p.l) && (p.buff[p.cursor] == ' ') {
			p.cursor++
		}

		return ts, syslogparser.ErrTimestampUnknownFormat
	}

	fixTimestampIfNeeded(&ts)

	p.cursor += tsFmtLen

	if (p.cursor < p.l) && (p.buff[p.cursor] == ' ') {
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

func (p *Parser) parseApp() (string, string, error) {
	var b byte
	var endOfTag, closeBracket, openBracket bool
	var app, pid []byte
	var foundApp, foundPid bool

	from := p.cursor

	for {
		b = p.buff[p.cursor]

		openBracket = (b == '[')
		closeBracket = (b == ']')
		endOfTag = (b == ':' || b == ' ')

		if openBracket {
			app = p.buff[from:p.cursor]
			from = p.cursor
			foundApp = true
		} else if closeBracket {
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
				pid = p.buff[from:p.cursor]
				foundPid = true
			}

			p.cursor++
			break
		}

		p.cursor++
	}

	if (p.cursor < p.l) && (p.buff[p.cursor] == ' ') {
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

	newTs := time.Date(y, ts.Month(), ts.Day(),
		ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location())

	*ts = newTs
}
