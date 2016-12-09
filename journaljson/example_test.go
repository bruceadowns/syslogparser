package journaljson_test

import (
	"testing"

	"github.com/bruceadowns/syslogparser/journaljson"
)

func TestExampleNewParser(t *testing.T) {
	b := `{ "__CURSOR" : "s=36f3322f7ed84791baa20c2184f7f04e;i=48ad7fb;b=9f01452e66564f7eaf6f642f67b1fb55;m=58d70c73fdf;t=5433a929f270a;x=e12f7a4c83217fca", "__REALTIME_TIMESTAMP" : "1481293730621194", "__MONOTONIC_TIMESTAMP" : "6105040633823", "_BOOT_ID" : "9f01452e66564f7eaf6f642f67b1fb55", "PRIORITY" : "3", "MESSAGE" : "I1209 14:28:50.590700       1 handlers.go:152] GET /api/v1/namespaces/dev/pods?labelSelector=deployment%3D1ae01ff6ff28a41eb0a6f7d5a2152315%2Cname%3Dhello-world: (17.937289ms) 200 [[healthcheck/v0.0.0 (linux/amd64) kubernetes/$Format] 10.121.79.211:57084]", "CONTAINER_ID" : "db7a3c4909b2", "CONTAINER_ID_FULL" : "db7a3c4909b2e4760d47992a546afb1b16290390dc693e102bcd2b13e51e621b", "CONTAINER_NAME" : "k8s_kube-apiserver.e138433c_kube-apiserver-ip-10-121-113-159.us-west-2.compute.internal_kube-system_b588e28c8548665083300bb409cacb39_278c7649", "_TRANSPORT" : "journal", "_PID" : "4202", "_UID" : "0", "_GID" : "0", "_COMM" : "docker", "_EXE" : "/usr/bin/docker", "_CMDLINE" : "docker daemon --host=fd:// --log-driver=journald --exec-opt native.cgroupdriver=systemd --bip=172.24.39.1/24 --mtu=8951 --ip-masq=false --selinux-enabled", "_CAP_EFFECTIVE" : "3fffffffff", "_SYSTEMD_CGROUP" : "/system.slice/docker.service", "_SYSTEMD_UNIT" : "docker.service", "_SYSTEMD_SLICE" : "system.slice", "_SELINUX_CONTEXT" : "system_u:system_r:kernel_t:s0", "_SOURCE_REALTIME_TIMESTAMP" : "1481293730590954", "_MACHINE_ID" : "c1d87446d7894c78b2d5d96bb55d5cad", "_HOSTNAME" : "ip-10-121-113-159.us-west-2.compute.internal" }`
	buff := []byte(b)

	p := journaljson.NewParser(buff)

	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}

	t.Log(p.Dump())
}
