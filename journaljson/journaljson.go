package journaljson

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bruceadowns/syslogparser"
)

// JSON holds mako structured json
type JSON struct {
	Cursor                  string `json:"__CURSOR,omitempty"`
	RealtimeTimestamp       string `json:"__REALTIME_TIMESTAMP,omitempty"`
	MonotonicTimestamp      string `json:"__MONOTONIC_TIMESTAMP,omitempty"`
	BootID                  string `json:"_BOOT_ID,omitempty"`
	Priority                string `json:"PRIORITY,omitempty"`
	Message                 string `json:"MESSAGE,omitempty"`
	ContainerID             string `json:"CONTAINER_ID,omitempty"`
	ContainerIDFull         string `json:"CONTAINER_ID_FULL,omitempty"`
	ContainerName           string `json:"CONTAINER_NAME,omitempty"`
	Transport               string `json:"_TRANSPORT,omitempty"`
	PID                     string `json:"_PID,omitempty"`
	UID                     string `json:"_UID,omitempty"`
	GID                     string `json:"_GID,omitempty"`
	Comm                    string `json:"_COMM,omitempty"`
	Exe                     string `json:"_EXE,omitempty"`
	CmdLine                 string `json:"_CMDLINE,omitempty"`
	CapEffective            string `json:"_CAP_EFFECTIVE,omitempty"`
	SystemdCGroup           string `json:"_SYSTEMD_CGROUP,omitempty"`
	SystemdUnit             string `json:"_SYSTEMD_UNIT,omitempty"`
	SystemDSlice            string `json:"_SYSTEMD_SLICE,omitempty"`
	SeLinuxContext          string `json:"_SELINUX_CONTEXT,omitempty"`
	SourceRealtimeTimestamp string `json:"_SOURCE_REALTIME_TIMESTAMP,omitempty"`
	MachineID               string `json:"_MACHINE_ID,omitempty"`
	HostName                string `json:"_HOSTNAME,omitempty"`
}

// Parser struct
type Parser struct {
	bb          *bytes.Buffer
	journalJSON JSON
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		bb: bytes.NewBuffer(buff),
	}
}

// Parse ...
func (p *Parser) Parse() error {
	err := json.NewDecoder(p.bb).Decode(&p.journalJSON)
	if err != nil {
		return err
	}

	if len(p.journalJSON.HostName) == 0 {
		return fmt.Errorf("Host name not found")
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	level := "INFO"
	switch p.journalJSON.Priority {
	case "0", "1", "2", "3":
		level = "ERROR"
	case "4":
		level = "WARN"
	case "7":
		level = "DEBUG"
	//case "5", "6":
	default:
		level = "INFO"
	}

	timestamp := p.journalJSON.SourceRealtimeTimestamp
	if len(timestamp) == 0 {
		timestamp = p.journalJSON.RealtimeTimestamp
	}
	if len(timestamp) == 16 {
		timestamp = timestamp[:len(timestamp)-3]
	}

	return syslogparser.LogParts{
		"hostname":        p.journalJSON.HostName,
		"logger_name":     p.journalJSON.Transport,
		"level":           level,
		"message":         p.journalJSON.Message,
		"service_name":    p.journalJSON.Exe,
		"service_version": p.journalJSON.PID,
		"timestamp":       timestamp,
	}
}
