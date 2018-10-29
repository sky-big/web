package clog

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
)

type LogContext struct {
	depth int
	tag   map[string]string
	buf   bytes.Buffer
	Schema
}

var logPool = &sync.Pool{New: func() interface{} { return new(LogContext) }}

func getLogContext(s Schema) *LogContext {
	lc := logPool.Get().(*LogContext)
	lc.Schema = s
	lc.Reset()
	return lc
}

func putLogContext(lc *LogContext) {
	logPool.Put(lc.Reset())
}

func (lc *LogContext) setDepth(depth int) *LogContext {
	lc.depth = depth
	return lc
}

func (lc *LogContext) SetTag(key, value string) *LogContext {
	if lc.tag == nil {
		lc.tag = make(map[string]string)
	}
	lc.tag[key] = value
	return lc
}

func (lc *LogContext) SetRequestID(requestId string) *LogContext {
	lc.SetTag(LF_ReqID, requestId)
	return lc
}

func (lc *LogContext) SetRequest(req *http.Request) *LogContext {
	if req == nil {
		return lc
	}
	lc.SetTag(LF_IP, req.RemoteAddr)
	lc.SetTag(LF_HostURL, req.Host)
	lc.SetTag(LF_Method, req.Method)
	lc.SetTag(LF_Path, req.URL.RequestURI())
	lc.SetTag(LF_ReqID, req.Header.Get("requestId"))
	lc.SetTag(LF_RequestContentLength, strconv.FormatInt(req.ContentLength, 10))

	return lc
}

func (lc *LogContext) Reset() *LogContext {
	lc.tag = nil
	lc.buf.Reset() // clean buf
	return lc
}

func (lc *LogContext) Debug(args ...interface{}) {
	lc.SetTag(LF_Level, DEBUG)
	lc.Schema.head(lc)
	logging.print(debugLog, lc.buf.String(), args...)
	putLogContext(lc)
}

func (lc *LogContext) Debugf(format string, args ...interface{}) {
	lc.SetTag(LF_Level, DEBUG)
	lc.Schema.head(lc)
	logging.printf(debugLog, lc.buf.String(), format, args...)
	putLogContext(lc)
}

func (lc *LogContext) Info(args ...interface{}) {
	lc.SetTag(LF_Level, INFO)
	lc.Schema.head(lc)
	logging.print(infoLog, lc.buf.String(), args...)
	putLogContext(lc)
}

func (lc *LogContext) Infof(format string, args ...interface{}) {
	lc.SetTag(LF_Level, INFO)
	lc.Schema.head(lc)
	logging.printf(infoLog, lc.buf.String(), format, args...)
	putLogContext(lc)
}

func (lc *LogContext) Warning(args ...interface{}) {
	lc.SetTag(LF_Level, WARNING)
	lc.Schema.head(lc)
	logging.print(warningLog, lc.buf.String(), args...)
	putLogContext(lc)
}

func (lc *LogContext) Warningf(format string, args ...interface{}) {
	lc.SetTag(LF_Level, WARNING)
	lc.Schema.head(lc)
	logging.printf(warningLog, lc.buf.String(), format, args...)
	putLogContext(lc)
}

func (lc *LogContext) Error(args ...interface{}) {
	lc.SetTag(LF_Level, ERROR)
	lc.Schema.head(lc)
	logging.print(errorLog, lc.buf.String(), args...)
	putLogContext(lc)
}

func (lc *LogContext) Errorf(format string, args ...interface{}) {
	lc.SetTag(LF_Level, ERROR)
	lc.Schema.head(lc)
	logging.printf(errorLog, lc.buf.String(), format, args...)
	putLogContext(lc)
}

func (lc *LogContext) Fatal(args ...interface{}) {
	lc.SetTag(LF_Level, FATAL)
	lc.Schema.head(lc)
	logging.print(fatalLog, lc.buf.String(), args...)
	putLogContext(lc)
}

func (lc *LogContext) Fatalf(format string, args ...interface{}) {
	lc.SetTag(LF_Level, FATAL)
	lc.Schema.head(lc)
	logging.printf(fatalLog, lc.buf.String(), format, args...)
	putLogContext(lc)
}
