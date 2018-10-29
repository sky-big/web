package http_helpers

import (
	. "common/marshaler"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//unmarshal helper
func UnmarshalStream(marshaler Marshaler, reader io.Reader, item interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = marshaler.Unmarshal(data, item)
	if err != nil {
		return fmt.Errorf("unmarshal data error(%s), data: %s", err.Error(), string(data))
	}
	return nil
}

//return proper marshaler by http header
//json as default
func makeMarshaler(req *http.Request) Marshaler {
	for _, v := range req.Header["Content-Type"] {
		if v == "application/x-yaml" {
			return NewYamlMarshaler()
		}
	}
	return NewJsonMarshaler()
}

func UnmarshalHttpBody(req *http.Request, item interface{}) error {
	return UnmarshalStream(makeMarshaler(req), req.Body, item)
}
