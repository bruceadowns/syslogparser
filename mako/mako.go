package mako

import (
	"bytes"
	"encoding/json"

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
}

// Parser struct
type Parser struct {
	bb       *bytes.Buffer
	MakoJSON JSON
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		bb: bytes.NewBuffer(buff),
	}
}

// Parse ...
func (p *Parser) Parse() error {
	err := json.NewDecoder(p.bb).Decode(&p.MakoJSON)
	if err != nil {
		return err
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	return syslogparser.LogParts{
		"logger_name":         p.MakoJSON.LoggerName,
		"level":               p.MakoJSON.Level,
		"level_value":         p.MakoJSON.LevelValue,
		"message":             p.MakoJSON.Message,
		"service_environment": p.MakoJSON.ServiceEnvironment,
		"service_name":        p.MakoJSON.ServiceName,
		"service_pipeline":    p.MakoJSON.ServicePipeline,
		"service_version":     p.MakoJSON.ServiceVersion,
		"thread_name":         p.MakoJSON.ThreadName,
		"@timestamp":          p.MakoJSON.Timestamp,
		"version":             p.MakoJSON.Version,
	}
}
