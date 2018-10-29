package clog

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

type severity int32 // sync/atomic int32

const (
	debugLog severity = iota
	infoLog
	warningLog
	errorLog
	fatalLog
	numSeverity = 5
)

const (
	DEBUG   = "DEBUG"
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
	FATAL   = "FATAL"
)

var severityName = []string{
	debugLog:   DEBUG,
	infoLog:    INFO,
	warningLog: WARNING,
	errorLog:   ERROR,
	fatalLog:   FATAL,
}

func severityByName(s string) (severity, bool) {
	s = strings.ToUpper(s)
	for i, name := range severityName {
		if name == s {
			return severity(i), true
		}
	}
	return 0, false
}

// output level
var outputSeverity = debugLog

func SetLogLevel(outputLevel string) {
	severity, ok := severityByName(outputLevel)
	if !ok {
		panic(fmt.Errorf("unknown severity name %s", outputLevel))
	}

	outputSeverity = severity
}

// get returns the value of the severity.
func (s *severity) get() severity {
	return severity(atomic.LoadInt32((*int32)(s)))
}

// set sets the value of the severity.
func (s *severity) set(val severity) {
	atomic.StoreInt32((*int32)(s), int32(val))
}

func (s *severity) Set(value string) error {
	var threshold severity
	// Is it a known name?
	if v, ok := severityByName(value); ok {
		threshold = v
	} else {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		threshold = severity(v)
	}
	s.set(threshold)
	return nil
}
