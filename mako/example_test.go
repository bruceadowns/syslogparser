package mako_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/mako"
)

func TestExampleNewParserSso(t *testing.T) {
	b := `{"@timestamp":"2016-11-23T00:24:23.288+00:00","@version":1,"message":"172.16.23.1 - - [23/Nov/2016:00:24:23 +0000] \"POST /api/v1/auth?jive_instance_id=75fb6e39-fb99-4cd7-bff5-87b5a690510a HTTP/1.1\" 200 61 \"-\" \"Jakarta Commons-HttpClient/3.1\" 4","logger_name":"http.request","thread_name":"dw-11194","level":"INFO","level_value":20000,"service_name":"ps-sso-telefonica-o2","service_environment":"k8s-prod-ams","service_pipeline":"main","service_version":"37c469d0cea5d11fb62002740bc3b30496e014cd","pod_name":"ps-sso-telefonica-o2-868218942-g48th"}`
	buff := []byte(b)

	p := mako.NewParser(buff, "foobar")

	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserSampleGolang(t *testing.T) {
	b := `{"@timestamp":"2016-11-11T11:34:40-08:00","service_name":"sample-golang","service_environment":"mako-dev-test","service_pipeline":"main","service_version":"15713b5b733e11069f68c2b78cbe5ad3e4a40abc","message":"&{GET / HTTP/1.1 1 1 map[User-Agent:[curl/7.49.1] Accept:[* /*]] 0x3cbc70 0 [] false 127.0.0.1:8080 map[] map[] <nil> map[] 127.0.0.1:64304 / <nil> <nil> <nil> 0xc4200ee4c0}","level":"INFO"}`
	buff := []byte(b)

	p := mako.NewParser(buff, "foobar")

	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserJcx(t *testing.T) {
	b := `{"@timestamp":"2016-11-22T23:27:00.180+00:00","@version":1,"message":"172.16.43.0 - developer [22/Nov/2016:23:27:00 +0000] \"GET /tasks/7645a092-6c37-4d7b-b2a7-832b5969cd4d HTTP/1.1\" 200 - \"https://cloud-jcx-api.ms-integ.svc.jivehosted.com/ui/\" \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/53.0.2785.143 Chrome/53.0.2785.143 Safari/537.36\" 67","logger_name":"http.request","thread_name":"dw-609","level":"INFO","level_value":20000,"jcx.inst.uri":"jcx-inst-m5vql0n7fbpy1tvrb8c8wj","service_name":"cloud-jcx-api","service_environment":"ms-integ","service_pipeline":"main","service_version":"d6b90788605888b6bc814d584cc529fcba9a300a"}`
	buff := []byte(b)

	p := mako.NewParser(buff, "foobar")

	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}

func TestExampleNewParserMituiCloudalytics(t *testing.T) {
	b := `{"service_name":"mitui-cloudalytics","service_environment":"ms-integ","service_pipeline":"main","service_version":"73600068e8d2648d3c1d33e5f1db6b0e039c3a90","k8s_pod_name":"mitui-cloudalytics-3757238288-j5lzg","hostname":"mitui-cloudalytics-3757238288-j5lzg","pid":16,"className":"adminServer","level":30,"message":"MITUI admin server listening on port 8081","timestamp":"2016-12-06T01:33:11.441Z","v":0}`
	buff := []byte(b)

	p := mako.NewParser(buff, "foobar")

	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
