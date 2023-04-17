package db

import (
	"fmt"
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/model/filter"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
	"time"
)

func TestDb(t *testing.T) {
	var (
		client *Db
		vp     *viper.Viper
		err    error
	)
	if vp, err = config.InitConfig(filepath.Join("..", "..", "configs"), "config"); err != nil {
		t.Error(err)
	}
	if client, err = New(
		vp.GetString("postgres.login"),
		GetEnv("POSTGRES_PASSWORD", "root"),
		vp.GetString("postgres.host"),
		vp.GetString("postgres.port"),
		vp.GetString("postgres.name"),
	); err != nil {
		t.Error(err)
	}

	filters := []*filter.Filter{
		{Filter: "filter1"},
		{Filter: "filter2"},
		{Filter: "filter3"},
	}
	filtrationRules := []*filter.FiltrationRule{
		{
			Filter:         filters[0],
			FilterField:    "field1",
			FilterFunction: "function1",
			FilterValue:    "value1",
			UpdatedAt:      time.Now().Add(time.Second * 5),
		},
		{
			FilterField:    "field2",
			FilterFunction: "function2",
			FilterValue:    "value2",
			UpdatedAt:      time.Now(),
		},
	}

	t.Run("create_filter", func(t2 *testing.T) {
		err := client.CreateFilter(filters[0])
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(fmt.Sprintf("filter created successfully: id = %d", filters[0].Id))
	})

	t.Run("create_filtration_rule", func(t2 *testing.T) {
		err := client.CreateFiltrationRule(filtrationRules[0])
		t2.Log(filtrationRules[0].Filter)
		if err != nil {
			t2.Error(err)
			return
		}
		t2.Log("filtration rule created")
	})
	fmt.Println(client)

	t.Run("get_latest_filtration_rule", func(t2 *testing.T) {
		fr, err := client.GetLatestFiltrationRule()
		if err != nil {
			t2.Error(err)
			return
		}
		t2.Log(fr)
	})
	fmt.Println(client)

	t.Run("update_filtration_rule", func(t2 *testing.T) {
		err := client.UpdateFiltrationRule(filtrationRules[0])
		if err != nil {
			t2.Error(err)
			return
		}
		t2.Log("Filtration rule updated")
	})

	//t.Run("get_filtration_rule", func(t2 *testing.T) {
	//	fr, err := client.GetFiltrationRule(1)
	//	if err != nil {
	//		t2.Error(err)
	//	}
	//	t2.Log(fr)
	//})
	//fmt.Println(client)
}
