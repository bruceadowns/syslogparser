package mako

import (
	"bytes"
	"encoding/json"
	"log"
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
	Timestamp          string `json:"@timestamp,omitempty"`
	Version            int    `json:"@version,omitempty"`
	V                  int    `json:"v,omitempty"`
}

// Parser struct
type Parser struct {
	bb       *bytes.Buffer
	hostname string
	MakoJSON JSON
}

// NewParser ...
func NewParser(buff []byte, hostname string) *Parser {
	return &Parser{
		bb:       bytes.NewBuffer(buff),
		hostname: hostname,
	}
}

// Parse ...
func (p *Parser) Parse() error {
	bb := bytes.NewBufferString(strings.Replace(p.bb.String(), ",\"level\":30,", ",\"level\":\"INFO\",", -1))

	err := json.NewDecoder(bb).Decode(&p.MakoJSON)
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
		"thread_name":         p.MakoJSON.ThreadName,
		"@timestamp":          timestamp,
		"version":             version,
	}
}
