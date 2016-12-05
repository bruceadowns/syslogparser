package syslogparser

import "strconv"

// Constants
const (
	priPartStart = '<'
	priPartEnd   = '>'
	NoVersion    = -1
)

// Err constants
var (
	ErrEOL                    = &ParserError{"End of log line"}
	ErrNoSpace                = &ParserError{"No space found"}
	ErrPriorityNoStart        = &ParserError{"No start char found for priority"}
	ErrPriorityEmpty          = &ParserError{"Priority field empty"}
	ErrPriorityNoEnd          = &ParserError{"No end char found for priority"}
	ErrPriorityTooShort       = &ParserError{"Priority field too short"}
	ErrPriorityTooLong        = &ParserError{"Priority field too long"}
	ErrPriorityNonDigit       = &ParserError{"Non digit found in priority"}
	ErrVersionNotFound        = &ParserError{"Cannot find version"}
	ErrTimestampUnknownFormat = &ParserError{"Timestamp format unknown"}
	ErrYearInvalid            = &ParserError{"Invalid year in timestamp"}
	ErrMonthInvalid           = &ParserError{"Invalid month in timestamp"}
	ErrDayInvalid             = &ParserError{"Invalid day in timestamp"}
	ErrHourInvalid            = &ParserError{"Invalid hour in timestamp"}
	ErrMinuteInvalid          = &ParserError{"Invalid minute in timestamp"}
	ErrSecondInvalid          = &ParserError{"Invalid second in timestamp"}
	ErrSecFracInvalid         = &ParserError{"Invalid fraction of second in timestamp"}
	ErrTimeZoneInvalid        = &ParserError{"Invalid time zone in timestamp"}
	ErrInvalidTimeFormat      = &ParserError{"Invalid time format"}
	ErrInvalidAppName         = &ParserError{"Invalid app name"}
	ErrInvalidProcID          = &ParserError{"Invalid proc ID"}
	ErrInvalidMsgID           = &ParserError{"Invalid msg ID"}
	ErrNoStructuredData       = &ParserError{"No structured data"}
)

// ParserError ...
type ParserError struct {
	ErrorString string
}

func (err *ParserError) Error() string {
	return err.ErrorString
}

// LogParser ...
type LogParser interface {
	Parse() error
	Dump() LogParts
}

// Priority ...
type Priority struct {
	P int
	F Facility
	S Severity
}

// Facility ...
type Facility struct {
	Value int
}

// Severity ...
type Severity struct {
	Value int
}

// LogParts ...
type LogParts map[string]interface{}

// ParsePriority per https://tools.ietf.org/html/rfc3164#section-4.1
func ParsePriority(buff []byte, cursor *int, l int) (Priority, error) {
	pri := newPriority(0)

	if l <= 0 {
		return pri, ErrPriorityEmpty
	}

	if buff[*cursor] != priPartStart {
		return pri, ErrPriorityNoStart
	}

	i := 1
	priDigit := 0

	for i < l {
		if i >= 5 {
			return pri, ErrPriorityTooLong
		}

		c := buff[i]

		if c == priPartEnd {
			if i == 1 {
				return pri, ErrPriorityTooShort
			}

			*cursor = i + 1
			return newPriority(priDigit), nil
		}

		if IsDigit(c) {
			v, e := strconv.Atoi(string(c))
			if e != nil {
				return pri, e
			}

			priDigit = (priDigit * 10) + v
		} else {
			return pri, ErrPriorityNonDigit
		}

		i++
	}

	return pri, ErrPriorityNoEnd
}

// ParseVersion per https://tools.ietf.org/html/rfc5424#section-6.2.2
func ParseVersion(buff []byte, cursor *int, l int) (int, error) {
	if *cursor >= l {
		return NoVersion, ErrVersionNotFound
	}

	c := buff[*cursor]
	*cursor++

	// XXX : not a version, not an error though as RFC 3164 does not support it
	if !IsDigit(c) {
		return NoVersion, nil
	}

	v, e := strconv.Atoi(string(c))
	if e != nil {
		*cursor--
		return NoVersion, e
	}

	return v, nil
}

// IsDigit ...
func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func newPriority(p int) Priority {
	// The Priority value is calculated by first multiplying the Facility
	// number by 8 and then adding the numerical value of the Severity.

	return Priority{
		P: p,
		F: Facility{Value: p / 8},
		S: Severity{Value: p % 8},
	}
}

// FindNextSpace ...
func FindNextSpace(buff []byte, from int, l int) (int, error) {
	var to int

	for to = from; to < l; to++ {
		if buff[to] == ' ' {
			to++
			return to, nil
		}
	}

	return 0, ErrNoSpace
}

// Parse2Digits ...
func Parse2Digits(buff []byte, cursor *int, l int, min int, max int, e error) (int, error) {
	digitLen := 2

	if *cursor+digitLen > l {
		return 0, ErrEOL
	}

	sub := string(buff[*cursor : *cursor+digitLen])

	*cursor += digitLen

	i, err := strconv.Atoi(sub)
	if err != nil {
		return 0, e
	}

	if i >= min && i <= max {
		return i, nil
	}

	return 0, e
}

// ParseHostname ...
func ParseHostname(buff []byte, cursor *int, l int) (string, error) {
	from := *cursor
	var to int

	for to = from; to < l; to++ {
		if buff[to] == ' ' {
			break
		}
	}

	hostname := buff[from:to]

	*cursor = to

	return string(hostname), nil
}
