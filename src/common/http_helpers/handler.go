package http_helpers

import (
	. "common/clog"
	"net/http"
	"runtime/debug"
	"time"
	"github.com/gorilla/mux"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func Wrap(f HandlerFunc) HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		// print request trace info
		start := time.Now()
		defer func() {
			if r := recover(); r != nil {
				Hlog.SetRequest(req).
					SetTag(LF_ResponseTime, time.Since(start).String()).
					Errorf("Recover in Request, stack: %s", string(debug.Stack()))
			}
		}()
		Hlog.SetRequest(req).
			Info("Request start")
		f(w, req)
		Hlog.SetRequest(req).
			SetTag(LF_ResponseTime, time.Since(start).String()).
			Info("Request done")
	}
}

// make handler func more easy to be tested
type ThinHandlerFunc func(*http.Request) (
	data interface{}, /*api-specified response result*/
	err int,
	message string /*reason if err not empty*/)

func WrapTHF(f ThinHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, errorType, msg := f(req)
		res := new(Response)
		res.Code = errorType
		if errorType == NoError {
			res.Data = data
		} else {
			res.Message = msg
			res.Code = errorType
			Hlog.SetRequest(req).
				Debugf("Request error(%s) message(%s)", errorType, msg)
		}
		DoResponse(w, res)
	}
}

func AddHandleFunc(r *mux.Router, method string, path string, f HandlerFunc) {
	r.Methods(method).Path(path).HandlerFunc(Wrap(f))
}

func AddThinHandleFunc(r *mux.Router, method string, path string, f ThinHandlerFunc) {
	AddHandleFunc(r, method, path, WrapTHF(f))
}