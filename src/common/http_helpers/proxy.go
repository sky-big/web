package http_helpers

import (
	"net/http"
	"net/http/httputil"
	"bytes"
	"log"
	. "common/clog"
)

func Forward(p *httputil.ReverseProxy, w http.ResponseWriter, req *http.Request){
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Lshortfile)
	p.ErrorLog = logger
	p.ServeHTTP(w, req)
	if buf.Len() != 0 {
		Blog.Warning(buf.String())
	}
}
