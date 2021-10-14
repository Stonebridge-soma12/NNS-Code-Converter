package Config

import (
	"fmt"
	"go.uber.org/config"
	"io/ioutil"
	"strings"
)

type Config struct {
	Host    string `yaml:"host"`
	Account string `yaml:"account"`
	Pw      string `yaml:"password"`
	VHost   string `yaml:"vhost"`
}

func GetConfig() (Config, error) {
	var cf Config

	// base
	baseCfg, err := ioutil.ReadFile("./config/config.base.yaml")
	if err != nil {
		return cf, err
	}
	base := strings.NewReader(string(baseCfg))

	cfg, err := ioutil.ReadFile("./config/config.dev.yaml")
	if err != nil {
		return cf, err
	}
	conf := strings.NewReader(string(cfg))

	yaml, err := config.NewYAML(config.Source(base), config.Source(conf))
	if err != nil {
		return cf, err
	}

	err = yaml.Get("Config").Populate(&cf)
	if err != nil {
		return cf, err
	}

	fmt.Println(cf)

	return cf, nil
}
