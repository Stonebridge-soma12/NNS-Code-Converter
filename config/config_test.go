package Config

import (
	"fmt"
	"testing"
)

func Test_GetConfig(t *testing.T) {
	conf, err := GetConfig()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(conf.Account)
	fmt.Println(conf.Host)
}