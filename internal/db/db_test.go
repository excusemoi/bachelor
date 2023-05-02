package db

import (
	"fmt"
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/model"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
	"time"
)

func TestDb(t *testing.T) {
	var (
		client = &Db[model.FiltrationRule]{}
		vp     *viper.Viper
		err    error
	)
	if vp, err = config.InitConfig(filepath.Join("..", "..", "configs"), "config"); err != nil {
		t.Error(err)
	}
	if client, err = client.Init(vp); err != nil {
		t.Error(err)
	}

	filtrationRules := []*model.FiltrationRule{
		{
			Field:     "field1",
			Function:  "function1",
			Value:     "value1",
			UpdatedAt: time.Now().Add(time.Second * 5),
		},
		{
			Field:     "field2",
			Function:  "function2",
			Value:     "value2",
			UpdatedAt: time.Now(),
		},
	}

	t.Run("create_filtration_rule", func(t2 *testing.T) {
		err := client.Create(filtrationRules[0])
		if err != nil {
			t2.Error(err)
			return
		}
		t2.Log("filtration rule created")
	})
	fmt.Println(client)

	t.Run("get_latest", func(t2 *testing.T) {
		fr := &model.FiltrationRule{}
		_, err = client.GetLatest(fr)
		if err != nil {
			t2.Error(err)
			return
		}
		t2.Log("Filtration rule.id:" + fr.Value)
	})

	//t.Run("get_filtration_rule", func(t2 *testing.T) {
	//	fr, err := client.GetByID(1)
	//	if err != nil {
	//		t2.Error(err)
	//	}
	//	t2.Log(fr)
	//})
	//fmt.Println(client)
}
