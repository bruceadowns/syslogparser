package rfc3164_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/rfc3164"
)

func TestExampleNewParser(t *testing.T) {
	b := "<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"
	t.Log(b)
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewUbuntuRsyslog(t *testing.T) {
	b := "<46>Dec  7 17:51:54 bdowns-virtual-machine rsyslogd-2359: action 'action 1' resumed (module 'builtin:omfwd') [v8.16.0 try http://www.rsyslog.com/e/2359 ]'"
	t.Log(b)
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewUbuntuSystemd(t *testing.T) {
	b := `<30>Dec  7 17:52:24 bdowns-virtual-machine systemd[1211]: Stopped target Timers.`
	t.Log(b)
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewOsx(t *testing.T) {
	b := ``
	t.Log(b)
	buff := []byte(b)

	p := rfc3164.NewParser(buff)
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
