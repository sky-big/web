package conf

import (
	"io/ioutil"

	. "common/clog"

	"github.com/go-yaml/yaml"
)

type Configure struct {
	Port      string    `yaml:"port"`
	HTMLPath  string    `yaml:"html_path"`
	LogConfig LogConfig `yaml:"log_config"`
}

var C *Configure

func LoadConfig(filename string) (*Configure, error) {
	c := new(Configure)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Init(filename string) error {
	var err error
	C, err = LoadConfig(filename)
	if err != nil {
		Blog.Errorf("matrix load config error : %s", err.Error())
		return err
	}

	// init log
	LogInit(C.LogConfig)

	return err
}
