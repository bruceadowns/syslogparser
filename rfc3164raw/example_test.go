package rfc3164raw_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/rfc3164raw"
)

func TestExampleNewParserPid(t *testing.T) {
	b := "Dec  1 00:03:16 ip-10-126-5-155 dhclient[2346]: bound to 10.126.5.155 -- renewal in 1721 seconds."
	t.Log(b)
	buff := []byte(b)

	p := rfc3164raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserBase(t *testing.T) {
	b := "Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"
	t.Log(b)
	buff := []byte(b)

	p := rfc3164raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserEc2(t *testing.T) {
	b := "Dec  1 00:03:17 ip-10-126-5-155 ec2net: [rewrite_aliases] Rewriting aliases of eth0"
	t.Log(b)
	buff := []byte(b)

	p := rfc3164raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
