package http_helpers

import (
	. "common/clog"
	"encoding/json"
	"net/http"
	"fmt"
)

const (
	NoError   = 0
	NoMessage = ""
)

/**
 * Unified response struct for http request
 * Will be marshaled to json and put in http body
 */
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // api-specified response result
}

var internalError = []byte("{\"error\": \"internal\"}")

//response helper
func DoResponse(w http.ResponseWriter, response *Response) {
	header := w.Header()
	header["Content-Type"] = []string{"application/json"}
	bin, err := json.Marshal(response)
	if err != nil {
		w.Write(internalError)
		Blog.Errorf("Marshal response error(%s) body(%+v) ", err.Error(), response)
	} else {
		w.Write(bin)
	}
}

type OpenAPIResponse struct {
	RequestId string      `json:"requestId"`
	Result    interface{} `json:"result"`
	Error     interface{} `json:"error"`
}

func DoOpenAPIResponse(w http.ResponseWriter, response *OpenAPIResponse) {
	header := w.Header()
	header["Content-Type"] = []string{"application/json"}
	bin, err := json.Marshal(response)
	if err != nil {
		w.Write(internalError)
		Blog.Errorf("Marshal response error(%s) body(%+v) ", err.Error(), response)
	} else {
		w.Write(bin)
	}
}

func MakeResponse(w http.ResponseWriter, req *http.Request, reCode int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args)
	response := Response{ Code: reCode, Message: message}
	ret, _ := json.Marshal(response)
	w.Write(ret)
}
