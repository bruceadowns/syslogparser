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
	MakoJSON JSON
}

// NewParser ...
func NewParser(buff []byte, hostname net.Addr) *Parser {
	return &Parser{
		bb:       bytes.NewBuffer(buff),
		hostname: hostname.String(),
	}
}

// global const in order to compile once
var reVersionStrung = regexp.MustCompile("\"version\":\"[0-9.]+\"")

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
	err := json.NewDecoder(preProcess(p.bb)).Decode(&p.MakoJSON)
	if err != nil {
		return err
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	levelValue := strconv.Itoa(p.MakoJSON.LevelValue)

	timestamp := "0"
	if ts, err := time.Parse(time.RFC3339, p.MakoJSON.Timestamp); err == nil {
		timestamp = syslogparser.Epoch(ts)
	} else {
		log.Printf("Error parsing timestamp: %s", err)
	}

	version := strconv.Itoa(p.MakoJSON.Version)

	return syslogparser.LogParts{
		"hostname":            p.hostname,
		"logger_name":         p.MakoJSON.LoggerName,
		"level":               p.MakoJSON.Level,
		"level_value":         levelValue,
		"message":             p.MakoJSON.Message,
		"service_environment": p.MakoJSON.ServiceEnvironment,
		"service_name":        p.MakoJSON.ServiceName,
		"service_pipeline":    p.MakoJSON.ServicePipeline,
		"service_version":     p.MakoJSON.ServiceVersion,
		"stack_trace":         p.MakoJSON.StackTrace,
		"thread_name":         p.MakoJSON.ThreadName,
		"timestamp":           timestamp,
		"version":             version,
	}
}