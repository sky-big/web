package clog

import (
	"testing"
)

func Test_GetTimestamp(t *testing.T) {
	filename := "jvessel-worker.ERROR.2018.01.02-18.12.13"
	ts, err := getTimestamp("jvessel-worker", "ERROR", filename)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(ts.String())
}

/*
cd /root/logs
touch jvessel-worker.ERROR.2018.11.11-11.11.11
touch jvessel-worker.ERROR.2018.11.11-11.11.12
touch jvessel-worker.ERROR.2018.11.11-11.11.13
touch jvessel-worker.ERROR.2018.11.11-11.11.14
touch jvessel-worker.ERROR.2018.11.11-11.11.09
touch jvessel-worker.ERROR.2018.11.11-11.11.05
touch jvessel-worker.ERROR.2018.11.11-11.11.07
touch jvessel-worker.ERROR.2018.11.11-11.11.01
*/
func Test_Delete(t *testing.T) {
	return
	rotate("jvessel-worker", "ERROR", "/root/logs")
}
