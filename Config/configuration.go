package Config

import (
	"fmt"
	"go.uber.org/config"
	"io/ioutil"
	"strings"
)

type Config struct {
	BaseURL string	`yaml:"base_url"`
	Port    string	`yaml:"port"`
}

func GetConfig() (Config, error) {
	var cf Config

	// base
	baseCfg, err := ioutil.ReadFile("/Config/config.base")
	if err != nil {
		return cf, err
	}
	base := strings.NewReader(string(baseCfg))

	cfg, err := ioutil.ReadFile("/Config/config.dev")
	if err != nil {
		return cf, err
	}
	conf := strings.NewReader(string(cfg))

	fmt.Println(conf)

	yaml, err := config.NewYAML(config.Source(base), config.Source(conf))
	if err != nil {
		return cf, err
	}

	err = yaml.Get("URL").Populate(&cf)
	if err != nil {
		return cf, err
	}

	return cf, nil
}