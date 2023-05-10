package kafka

import (
	"errors"
	"github.com/bachelor/internal/config"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestKafka(t *testing.T) {
	var (
		c   = &Client{}
		vp  = &viper.Viper{}
		err error
	)

	if vp, err = config.InitConfig(filepath.Join("..", "components", "source", "configs"), "config-local"); err != nil {
		t.Fatal(err)
	}

	c.Init(vp)

	t.Log("kafka handling events")

	c.HandleEvents(func(bytes []byte) ([]byte, error) {
		s := string(bytes)
		t.Log(s)
		if s == "error" {
			return nil, errors.New("")
		}
		return bytes, nil
	})
}
