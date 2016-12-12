package mako

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bruceadowns/syslogparser"
)

// JSON holds mako structured json
type JSON struct {
	LoggerName         string `json:"logger_name,omitempty"`
	Message            string `json:"message,omitempty"`
	Level              string `json:"level,omitempty"`
	LevelValue         int    `json:"level_value,omitempty"`
	ServiceEnvironment string `json:"service_environment,omitempty"`
	ServiceName        string `json:"service_name,omitempty"`
	ServicePipeline    string `json:"service_pipeline,omitempty"`
	ServiceVersion     string `json:"service_version,omitempty"`
	ThreadName         string `json:"thread_name,omitempty"`
	Timestamp          string `json:"timestamp,omitempty"`
	Version            int    `json:"version,omitempty"`
	StackTrace         string `json:"stack_trace,omitempty"`
}

// Parser struct
type Parser struct {
	bb       *bytes.Buffer
	hostname string
	makoJSON JSON
}

// NewParser ...
func NewParser(buff []byte, hostname net.Addr) *Parser {
	return &Parser{
		bb:       bytes.NewBuffer(buff),
		hostname: hostname.String(),
	}
}

// global const in order to compile once
var reVersionStrung = regexp.MustCompile("\"version\":\"[^\"]+\"")

func preProcess(in *bytes.Buffer) io.Reader {
	replacer := strings.NewReplacer(
		"\"level\":10,", "\"level\":\"TRACE\",",
		"\"level\":20,", "\"level\":\"DEBUG\",",
		"\"level\":30,", "\"level\":\"INFO\",",
		"\"level\":40,", "\"level\":\"WARN\",",
		"\"level\":50,", "\"level\":\"ERROR\",",
		"\"level\":60,", "\"level\":\"ERROR\",",
		"\"@timestamp\"", "\"timestamp\"",
		"\"@version\"", "\"version\"")

	out := replacer.Replace(in.String())
	out = reVersionStrung.ReplaceAllString(out, "\"version\":0")
	return bytes.NewBufferString(out)
}

// Parse ...
func (p *Parser) Parse() error {
	err := json.NewDecoder(preProcess(p.bb)).Decode(&p.makoJSON)
	if err != nil {
		return err
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	timestamp := "0"
	if ts, err := time.Parse(time.RFC3339, p.makoJSON.Timestamp); err == nil {
		timestamp = syslogparser.Epoch(ts)
	} else {
		log.Printf("Error parsing timestamp: %s", err)
	}

	return syslogparser.LogParts{
		"hostname":            p.hostname,
		"logger_name":         p.makoJSON.LoggerName,
		"level":               p.makoJSON.Level,
		"level_value":         strconv.Itoa(p.makoJSON.LevelValue),
		"message":             p.makoJSON.Message,
		"service_environment": p.makoJSON.ServiceEnvironment,
		"service_name":        p.makoJSON.ServiceName,
		"service_pipeline":    p.makoJSON.ServicePipeline,
		"service_version":     p.makoJSON.ServiceVersion,
		"stack_trace":         p.makoJSON.StackTrace,
		"thread_name":         p.makoJSON.ThreadName,
		"timestamp":           timestamp,
		"version":             strconv.Itoa(p.makoJSON.Version),
	}
}
