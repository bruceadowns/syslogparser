package journalmako_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/journalmako"
)

func TestExampleNewParser(t *testing.T) {
	b := `{ "__CURSOR" : "s=2daffcb205ee48f3addc3a7496820091;i=2d7a3aa;b=9bd8d0d20a8746be9fbb237d52dfae8c;m=12dcd6b47e2;t=5432a2749c670;x=53d7f81f7edfe28a", "__REALTIME_TIMESTAMP" : "1481223210387056", "__MONOTONIC_TIMESTAMP" : "1296231516130", "_BOOT_ID" : "9bd8d0d20a8746be9fbb237d52dfae8c", "_TRANSPORT" : "journal", "_PID" : "2592", "_UID" : "0", "_GID" : "0", "_COMM" : "docker", "_EXE" : "/usr/bin/docker", "_CMDLINE" : "docker daemon --host=fd:// --storage-driver=overlay --iptables=false --log-driver=journald --exec-opt native.cgroupdriver=systemd --bip=172.16.87.1/24 --mtu=1450 --ip-masq=false --selinux-enabled", "_CAP_EFFECTIVE" : "3fffffffff", "_SYSTEMD_CGROUP" : "/system.slice/docker.service", "_SYSTEMD_UNIT" : "docker.service", "_SYSTEMD_SLICE" : "system.slice", "_SELINUX_CONTEXT" : "system_u:system_r:kernel_t:s0", "_MACHINE_ID" : "5af5635132b1470f955156db2fa59ba2", "_HOSTNAME" : "iseb00240", "PRIORITY" : "6", "MESSAGE" : "{\"@timestamp\":\"2016-12-08T18:53:30.386+00:00\",\"@version\":1,\"message\":\"172.16.87.1 - - [08/Dec/2016:18:53:30 +0000] \\\"GET /ping HTTP/1.1\\\" 200 - \\\"-\\\" \\\"Go-http-client/1.1\\\" 1\",\"logger_name\":\"http.request\",\"thread_name\":\"dw-admin-19495\",\"level\":\"INFO\",\"level_value\":20000,\"service_name\":\"ps-addon-telefonica-head\",\"service_environment\":\"k8s-integ\",\"service_pipeline\":\"main\",\"service_version\":\"a4c52903f4af72ff62619dd574d2458f1e83c4ef\",\"pod_name\":\"ps-addon-telefonica-head-848342042-ueonv\"}", "CONTAINER_ID" : "8af4c7bf13c0", "CONTAINER_ID_FULL" : "8af4c7bf13c07009ce052e7aeed6332673945211adf59570670bc512df56645f", "CONTAINER_NAME" : "k8s_ps-addon-telefonica-head.44c9f80f_ps-addon-telefonica-head-848342042-ueonv_k8s-integ_f9dee9ca-b8a1-11e6-bfad-005056a74490_67f13527", "_SOURCE_REALTIME_TIMESTAMP" : "1481223210386785" }`
	t.Log(b)
	buff := []byte(b)

	p := journalmako.NewParser(buff)

	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
