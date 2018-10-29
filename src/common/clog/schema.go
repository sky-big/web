package clog

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Schema interface {
	// render log head in log context buf
	head(lc *LogContext)

	SetTag(key, value string) *LogContext
	SetRequest(req *http.Request) *LogContext

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type field struct {
	key   string // field key
	value string // field default value
}

type schema struct {
	fields []*field
}

func (s *schema) head(lc *LogContext) {
	_, file, line, ok := runtime.Caller(lc.depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	lc.SetTag(LF_ProcName, program)
	lc.SetTag(LF_Host, host)
	lc.SetTag(LF_FileLine, fmt.Sprintf("%s:%d", file, line))
	lc.SetTag(LF_Timestamp, time.Now().UTC().Add(8*time.Hour).Format("2006-01-02 15:04:05.999999"))
	for _, field := range s.fields {
		v, ok := lc.tag[field.key]
		if ok {
			lc.buf.WriteString(v) // set by user
		} else {
			lc.buf.WriteString(field.value) // default value
		}
		lc.buf.WriteString(" | ")
	}
}

func (s *schema) SetTag(key, value string) *LogContext {
	return getLogContext(s).setDepth(2).SetTag(key, value)
}

func (s *schema) SetRequestID(requestId string) *LogContext {
	return getLogContext(s).setDepth(2).SetRequestID(requestId)
}

func (s *schema) SetRequest(req *http.Request) *LogContext {
	return getLogContext(s).setDepth(2).SetRequest(req)
}

func (s *schema) Debug(args ...interface{}) {
	getLogContext(s).setDepth(3).Debug(args...)
}

func (s *schema) Debugf(format string, args ...interface{}) {
	getLogContext(s).setDepth(3).Debugf(format, args...)
}

func (s *schema) Info(args ...interface{}) {
	getLogContext(s).setDepth(3).Info(args...)
}

func (s *schema) Infof(format string, args ...interface{}) {
	getLogContext(s).setDepth(3).Infof(format, args...)
}

func (s *schema) Warning(args ...interface{}) {
	getLogContext(s).setDepth(3).Warning(args...)
}

func (s *schema) Warningf(format string, args ...interface{}) {
	getLogContext(s).setDepth(3).Warningf(format, args...)
}

func (s *schema) Error(args ...interface{}) {
	getLogContext(s).setDepth(3).Error(args...)
}

func (s *schema) Errorf(format string, args ...interface{}) {
	getLogContext(s).setDepth(3).Errorf(format, args...)
}

func (s *schema) Fatal(args ...interface{}) {
	getLogContext(s).setDepth(3).Fatal(args...)
}

func (s *schema) Fatalf(format string, args ...interface{}) {
	getLogContext(s).setDepth(3).Fatalf(format, args...)
}

// LF: log field
const (
	LF_Schema    = "schema"
	LF_Version   = "version"
	LF_ProcName  = "proc_name"
	LF_Host      = "host"
	LF_Timestamp = "timestamp"
	LF_Level     = "level"
	LF_FileLine  = "code_file_line"

	LF_SpaceID      = "space_id"
	LF_ReqID        = "request_id"
	LF_ResourceName = "resource_name"
	LF_ResourceKind = "resource_kind"

	LF_IP                   = "ip"
	LF_HostURL              = "host_url"
	LF_Method               = "method"
	LF_Path                 = "path"
	LF_RequestContentLength = "request_content_length"
	LF_ResponseTime         = "response_time"
)

// base log
var Blog = &schema{
	fields: []*field{
		// basic
		&field{key: LF_Schema, value: "JB"},
		&field{key: LF_Version, value: "v1"},
		&field{key: LF_ProcName},
		&field{key: LF_Host},
		&field{key: LF_Timestamp},
		&field{key: LF_Level},
		&field{key: LF_FileLine},
		&field{key: LF_ReqID},
	},
}

// event log
var Elog = &schema{
	fields: []*field{
		// basic
		&field{key: LF_Schema, value: "JE"},
		&field{key: LF_Version, value: "v1"},
		&field{key: LF_ProcName},
		&field{key: LF_Host},
		&field{key: LF_Timestamp},
		&field{key: LF_Level},
		&field{key: LF_FileLine},

		// option
		&field{key: LF_SpaceID},
		&field{key: LF_ResourceName},
		&field{key: LF_ResourceKind},
	},
}

// http log
var Hlog = &schema{
	fields: []*field{
		// basic
		&field{key: LF_Schema, value: "JH"},
		&field{key: LF_Version, value: "v1"},
		&field{key: LF_ProcName},
		&field{key: LF_Host},
		&field{key: LF_Timestamp},
		&field{key: LF_Level},
		&field{key: LF_FileLine},
		&field{key: LF_ReqID},

		// option
		&field{key: LF_IP},
		&field{key: LF_HostURL},
		&field{key: LF_Method},
		&field{key: LF_Path},
		&field{key: LF_RequestContentLength},
		&field{key: LF_ResponseTime},
	},
}
