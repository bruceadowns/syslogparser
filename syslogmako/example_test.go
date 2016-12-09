package syslogmako_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/syslogmako"
)

func TestExampleNewParser(t *testing.T) {
	b := `Dec 04 19:50:08 pseb00284 docker[2607]: {"@timestamp":"2016-12-04T19:50:08.221+00:00","@version":1,"message":"172.16.28.0 - - [04/Dec/2016:19:50:08 +0000] \"POST /api/v1/auth?jive_instance_id=75fb6e39-fb99-4cd7-bff5-87b5a690510a HTTP/1.1\" 200 61 \"-\" \"Jakarta Commons-HttpClient/3.1\" 3","logger_name":"http.request","thread_name":"dw-43354","level":"INFO","level_value":20000,"service_name":"ps-sso-telefonica-o2","service_environment":"k8s-prod-ams","service_pipeline":"main","service_version":"37c469d0cea5d11fb62002740bc3b30496e014cd","pod_name":"ps-sso-telefonica-o2-868218942-hc6s1"}`
	t.Log(b)
	buff := []byte(b)

	p := syslogmako.NewParser(buff)
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
