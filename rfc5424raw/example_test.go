package rfc5424raw_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/rfc5424raw"
)

func TestExampleNewParserRsyslog(t *testing.T) {
	b := `2016-11-27T03:34:01.968413-08:00 soa-prime-data1 rsyslogd: [origin software="rsyslogd" swVersion="5.8.10" x-pid="2169" x-info="http://www.rsyslog.com"] rsyslogd was HUPed`
	t.Log(b)
	buff := []byte(b)

	p := rfc5424raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserDd(t *testing.T) {
	b := `2016-11-30T16:08:49.538629-08:00 soa-prime-data1 dd.collector[2832]: INFO (collector.py:379): Finished run #1169220. Collection time: 4.09s. Emit time: 0.01s`
	t.Log(b)
	buff := []byte(b)

	p := rfc5424raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserGmetad(t *testing.T) {
	b := `2016-11-30T16:09:19.518625-08:00 soa-prime-data1 /usr/sbin/gmetad[2263]: data_thread() got no answer from any [my cluster] datasource`
	t.Log(b)
	buff := []byte(b)

	p := rfc5424raw.NewParser(buff)
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
