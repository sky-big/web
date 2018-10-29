package clog

import "fmt"

type LogConfig struct {
	Level  string `yaml:"level"`
	LogDir string `yaml:"dir"`
	// log to standard error instead of files
	ToStderr bool `yaml:"tostderr"`
	// log to standard error as well as files
	AlsoToStderr bool `yaml:"alsotostderr"`
	// logs at or above this threshold go to stderr
	StderrThreshold string `yaml:"stderrthreshold"`
}

func LogInit(c LogConfig) {
	if c.Level != "" {
		SetLogLevel(c.Level)
	}

	if c.LogDir != "" {
		SetLogDir(c.LogDir)
	}

	SetToStderr(c.ToStderr)

	SetAlsoToStderr(c.AlsoToStderr)

	if c.StderrThreshold != "" {
		SetStderrThreshold(c.StderrThreshold)
	}
}

// log to standard error instead of files
func SetToStderr(is bool) {
	logging.toStderr = is
}

// log to standard error as well as files
func SetAlsoToStderr(is bool) {
	logging.alsoToStderr = is
}

// logs at or above this threshold go to stderr
func SetStderrThreshold(s string) {
	level, ok := severityByName(s)
	if !ok {
		panic(fmt.Errorf("unknown severity name %s", s))
	}
	logging.stderrThreshold = level
}

func init() {
	logging.toStderr = true
	logging.alsoToStderr = false
	logging.stderrThreshold = errorLog
	go logging.flushDaemon()
}
