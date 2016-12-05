package rfc5424raw

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bruceadowns/syslogparser"
)

// Constants
const (
	NILVALUE = '-'
)

// Parser struct
type Parser struct {
	buff           []byte
	cursor         int
	l              int
	header         header
	structuredData string
	message        string
}

type header struct {
	timestamp time.Time
	hostname  string
	appName   string
	procID    string
}

type partialTime struct {
	hour    int
	minute  int
	seconds int
	secFrac float64
}

type fullTime struct {
	pt  partialTime
	loc *time.Location
}

type fullDate struct {
	year  int
	month int
	day   int
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		buff:   buff,
		cursor: 0,
		l:      len(buff),
	}
}

// Location ...
func (p *Parser) Location(location *time.Location) {
	// Ignore as RFC5424 syslog always has a timezone
}

// Parse ...
func (p *Parser) Parse() error {
	hdr, err := p.parseHeader()
	if err != nil {
		return err
	}

	p.header = hdr

	sd, err := p.parseStructuredData()
	if err == nil {
		p.structuredData = sd
		p.cursor++
	} else if err == syslogparser.ErrNoStructuredData {
		p.cursor++
	} else {
		return err
	}

	if p.cursor < p.l {
		p.message = string(p.buff[p.cursor:])
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	return syslogparser.LogParts{
		"timestamp":       p.header.timestamp,
		"hostname":        p.header.hostname,
		"app_name":        p.header.appName,
		"proc_id":         p.header.procID,
		"structured_data": p.structuredData,
		"message":         p.message,
	}
}

// HEADER = TIMESTAMP SP HOSTNAME SP APP-NAME[PROCID]:
func (p *Parser) parseHeader() (header, error) {
	hdr := header{}

	ts, err := p.parseTimestamp()
	if err != nil {
		return hdr, err
	}

	hdr.timestamp = ts
	p.cursor++

	host, err := p.parseHostname()
	if err != nil {
		return hdr, err
	}

	hdr.hostname = host
	p.cursor++

	appName, procID, err := p.parseApp()
	if err != nil {
		return hdr, err
	}

	hdr.appName = appName
	hdr.procID = procID

	return hdr, nil
}

// https://tools.ietf.org/html/rfc5424#section-6.2.3
func (p *Parser) parseTimestamp() (time.Time, error) {
	var ts time.Time

	if p.buff[p.cursor] == NILVALUE {
		p.cursor++
		return ts, nil
	}

	fd, err := parseFullDate(p.buff, &p.cursor, p.l)
	if err != nil {
		return ts, err
	}

	if p.buff[p.cursor] != 'T' {
		return ts, syslogparser.ErrInvalidTimeFormat
	}

	p.cursor++

	ft, err := parseFullTime(p.buff, &p.cursor, p.l)
	if err != nil {
		return ts, syslogparser.ErrTimestampUnknownFormat
	}

	nSec, err := toNSec(ft.pt.secFrac)
	if err != nil {
		return ts, err
	}

	ts = time.Date(
		fd.year,
		time.Month(fd.month),
		fd.day,
		ft.pt.hour,
		ft.pt.minute,
		ft.pt.seconds,
		nSec,
		ft.loc,
	)

	return ts, nil
}

// HOSTNAME = NILVALUE / 1*255PRINTUSASCII
func (p *Parser) parseHostname() (string, error) {
	return syslogparser.ParseHostname(p.buff, &p.cursor, p.l)
}

// app[pid]:
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

func (p *Parser) parseStructuredData() (string, error) {
	return parseStructuredData(p.buff, &p.cursor, p.l)
}

// ----------------------------------------------
// https://tools.ietf.org/html/rfc5424#section-6
// ----------------------------------------------

// FULL-DATE : DATE-FULLYEAR "-" DATE-MONTH "-" DATE-MDAY
func parseFullDate(buff []byte, cursor *int, l int) (fullDate, error) {
	var fd fullDate

	year, err := parseYear(buff, cursor, l)
	if err != nil {
		return fd, err
	}

	if buff[*cursor] != '-' {
		return fd, syslogparser.ErrTimestampUnknownFormat
	}

	*cursor++

	month, err := parseMonth(buff, cursor, l)
	if err != nil {
		return fd, err
	}

	if buff[*cursor] != '-' {
		return fd, syslogparser.ErrTimestampUnknownFormat
	}

	*cursor++

	day, err := parseDay(buff, cursor, l)
	if err != nil {
		return fd, err
	}

	fd = fullDate{
		year:  year,
		month: month,
		day:   day,
	}

	return fd, nil
}

// DATE-FULLYEAR   = 4DIGIT
func parseYear(buff []byte, cursor *int, l int) (int, error) {
	yearLen := 4

	if *cursor+yearLen > l {
		return 0, syslogparser.ErrEOL
	}

	// XXX : we do not check for a valid year (ie. 1999, 2013 etc)
	// XXX : we only checks the format is correct
	sub := string(buff[*cursor : *cursor+yearLen])

	*cursor += yearLen

	year, err := strconv.Atoi(sub)
	if err != nil {
		return 0, syslogparser.ErrYearInvalid
	}

	return year, nil
}

// DATE-MONTH = 2DIGIT  ; 01-12
func parseMonth(buff []byte, cursor *int, l int) (int, error) {
	return syslogparser.Parse2Digits(buff, cursor, l, 1, 12, syslogparser.ErrMonthInvalid)
}

// DATE-MDAY = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on month/year
func parseDay(buff []byte, cursor *int, l int) (int, error) {
	// XXX : this is a relaxed constraint
	// XXX : we do not check if valid regarding February or leap years
	// XXX : we only checks that day is in range [01 -> 31]
	// XXX : in other words this function will not rant if you provide Feb 31th
	return syslogparser.Parse2Digits(buff, cursor, l, 1, 31, syslogparser.ErrDayInvalid)
}

// FULL-TIME = PARTIAL-TIME TIME-OFFSET
func parseFullTime(buff []byte, cursor *int, l int) (fullTime, error) {
	var loc = new(time.Location)
	var ft fullTime

	pt, err := parsePartialTime(buff, cursor, l)
	if err != nil {
		return ft, err
	}

	loc, err = parseTimeOffset(buff, cursor, l)
	if err != nil {
		return ft, err
	}

	ft = fullTime{
		pt:  pt,
		loc: loc,
	}

	return ft, nil
}

// PARTIAL-TIME = TIME-HOUR ":" TIME-MINUTE ":" TIME-SECOND[TIME-SECFRAC]
func parsePartialTime(buff []byte, cursor *int, l int) (partialTime, error) {
	var pt partialTime

	hour, minute, err := getHourMinute(buff, cursor, l)
	if err != nil {
		return pt, err
	}

	if buff[*cursor] != ':' {
		return pt, syslogparser.ErrInvalidTimeFormat
	}

	*cursor++

	// ----

	seconds, err := parseSecond(buff, cursor, l)
	if err != nil {
		return pt, err
	}

	pt = partialTime{
		hour:    hour,
		minute:  minute,
		seconds: seconds,
	}

	// ----

	if buff[*cursor] != '.' {
		return pt, nil
	}

	*cursor++

	secFrac, err := parseSecFrac(buff, cursor, l)
	if err != nil {
		return pt, nil
	}
	pt.secFrac = secFrac

	return pt, nil
}

// TIME-HOUR = 2DIGIT  ; 00-23
func parseHour(buff []byte, cursor *int, l int) (int, error) {
	return syslogparser.Parse2Digits(buff, cursor, l, 0, 23, syslogparser.ErrHourInvalid)
}

// TIME-MINUTE = 2DIGIT  ; 00-59
func parseMinute(buff []byte, cursor *int, l int) (int, error) {
	return syslogparser.Parse2Digits(buff, cursor, l, 0, 59, syslogparser.ErrMinuteInvalid)
}

// TIME-SECOND = 2DIGIT  ; 00-59
func parseSecond(buff []byte, cursor *int, l int) (int, error) {
	return syslogparser.Parse2Digits(buff, cursor, l, 0, 59, syslogparser.ErrSecondInvalid)
}

// TIME-SECFRAC = "." 1*6DIGIT
func parseSecFrac(buff []byte, cursor *int, l int) (float64, error) {
	maxDigitLen := 6

	max := *cursor + maxDigitLen
	from := *cursor
	to := from

	for to = from; to < max; to++ {
		if to >= l {
			break
		}

		c := buff[to]
		if !syslogparser.IsDigit(c) {
			break
		}
	}

	sub := string(buff[from:to])
	if len(sub) == 0 {
		return 0, syslogparser.ErrSecFracInvalid
	}

	secFrac, err := strconv.ParseFloat("0."+sub, 64)
	*cursor = to
	if err != nil {
		return 0, syslogparser.ErrSecFracInvalid
	}

	return secFrac, nil
}

// TIME-OFFSET = "Z" / TIME-NUMOFFSET
func parseTimeOffset(buff []byte, cursor *int, l int) (*time.Location, error) {

	if buff[*cursor] == 'Z' {
		*cursor++
		return time.UTC, nil
	}

	return parseNumericalTimeOffset(buff, cursor, l)
}

// TIME-NUMOFFSET  = ("+" / "-") TIME-HOUR ":" TIME-MINUTE
func parseNumericalTimeOffset(buff []byte, cursor *int, l int) (*time.Location, error) {
	var loc = new(time.Location)

	sign := buff[*cursor]

	if (sign != '+') && (sign != '-') {
		return loc, syslogparser.ErrTimeZoneInvalid
	}

	*cursor++

	hour, minute, err := getHourMinute(buff, cursor, l)
	if err != nil {
		return loc, err
	}

	tzStr := fmt.Sprintf("%s%02d:%02d", string(sign), hour, minute)
	tmpTs, err := time.Parse("-07:00", tzStr)
	if err != nil {
		return loc, err
	}

	return tmpTs.Location(), nil
}

func getHourMinute(buff []byte, cursor *int, l int) (int, int, error) {
	hour, err := parseHour(buff, cursor, l)
	if err != nil {
		return 0, 0, err
	}

	if buff[*cursor] != ':' {
		return 0, 0, syslogparser.ErrInvalidTimeFormat
	}

	*cursor++

	minute, err := parseMinute(buff, cursor, l)
	if err != nil {
		return 0, 0, err
	}

	return hour, minute, nil
}

func toNSec(sec float64) (int, error) {
	_, frac := math.Modf(sec)
	fracStr := strconv.FormatFloat(frac, 'f', 9, 64)
	fracInt, err := strconv.Atoi(fracStr[2:])
	if err != nil {
		return 0, err
	}

	return fracInt, nil
}

// ------------------------------------------------
// https://tools.ietf.org/html/rfc5424#section-6.3
// ------------------------------------------------

func parseStructuredData(buff []byte, cursor *int, l int) (string, error) {
	var sdData string
	var found bool

	if buff[*cursor] == NILVALUE {
		*cursor++
		return "-", nil
	}

	if buff[*cursor] != '[' {
		return sdData, syslogparser.ErrNoStructuredData
	}

	from := *cursor
	to := from

	for to = from; to < l; to++ {
		if found {
			break
		}

		b := buff[to]

		if b == ']' {
			switch t := to + 1; {
			case t == l:
				found = true
			case t <= l && buff[t] == ' ':
				found = true
			}
		}
	}

	if found {
		*cursor = to
		return string(buff[from:to]), nil
	}

	return sdData, syslogparser.ErrNoStructuredData
}
