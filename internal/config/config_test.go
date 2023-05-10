package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		vp  = &viper.Viper{}
		err error
	)

	if vp, err = InitConfig(filepath.Join("..", "..", "configs"), "config-local"); err != nil {
		t.Fatal(err)
	}

	if err = vp.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	t.Logf("\"%v\"", nil == vp.Get("kafka.123"))

}
