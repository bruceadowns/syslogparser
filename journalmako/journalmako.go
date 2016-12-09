package journalmako

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/bruceadowns/syslogparser"
	"github.com/bruceadowns/syslogparser/journaljson"
	"github.com/bruceadowns/syslogparser/mako"
)

// Parser struct
type Parser struct {
	bb          *bytes.Buffer
	journalJSON journaljson.JSON
	makoJSON    mako.JSON
}

// NewParser ...
func NewParser(buff []byte) *Parser {
	return &Parser{
		bb: bytes.NewBuffer(buff),
	}
}

// global const in order to compile once
var reVersionStrung = regexp.MustCompile("\"version\":\"[0-9.]+\"")

func preProcess(in string) io.Reader {
	replacer := strings.NewReplacer(
		"\"level\":10,", "\"level\":\"TRACE\",",
		"\"level\":20,", "\"level\":\"DEBUG\",",
		"\"level\":30,", "\"level\":\"INFO\",",
		"\"level\":40,", "\"level\":\"WARN\",",
		"\"level\":50,", "\"level\":\"ERROR\",",
		"\"level\":60,", "\"level\":\"ERROR\",",
		"\"@timestamp\"", "\"timestamp\"",
		"\"@version\"", "\"version\"")

	out := replacer.Replace(in)
	out = reVersionStrung.ReplaceAllString(out, "\"version\":0")

	return bytes.NewBufferString(out)
}

// Parse ...
func (p *Parser) Parse() error {
	if err := json.NewDecoder(p.bb).Decode(&p.journalJSON); err != nil {
		return err
	}

	if err := json.NewDecoder(preProcess(p.journalJSON.Message)).Decode(&p.makoJSON); err != nil {
		return err
	}

	return nil
}

// Dump ...
func (p *Parser) Dump() syslogparser.LogParts {
	timestamp := p.journalJSON.SourceRealtimeTimestamp
	if len(timestamp) == 0 {
		timestamp = p.journalJSON.RealtimeTimestamp
	}
	if len(timestamp) == 16 {
		timestamp = timestamp[:len(timestamp)-3]
	}

	return syslogparser.LogParts{
		"hostname":            p.journalJSON.HostName,
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
