package main

import (
	. "common/clog"
)

func main() {
	LogInit(LogConfig{
		LogDir:       "/tmp/jvessel-test/",
		AlsoToStderr: true,
	})

	Elog.SetTag(LF_SpaceID, "sadfasf").
		Info("hi, this is info log")

	Elog.SetTag(LF_SpaceID, "sadfasf").
		SetResource(nil). //auto set: LF_SpaceID, LF_ServiceVersion, LF_ServiceName, LF_RequestID
		Debug("hi, this is info log")

	Elog.Info("hi", "h")
}
