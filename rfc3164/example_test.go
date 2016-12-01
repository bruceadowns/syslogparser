package rfc3164_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/rfc3164"
)

func TestExampleNewParser(t *testing.T) {
	b := "<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
