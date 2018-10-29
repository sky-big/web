// wrap json/yaml marshaler into unified interface
package marshaler

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type jsonMarshaler struct {
}

func NewJsonMarshaler() Marshaler {
	return new(jsonMarshaler)
}

func (t *jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (t *jsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type yamlMarshaler struct {
}

func NewYamlMarshaler() Marshaler {
	return new(yamlMarshaler)
}

func (t *yamlMarshaler) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (t *yamlMarshaler) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
