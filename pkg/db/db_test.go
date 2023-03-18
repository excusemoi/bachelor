package db

import (
	"fmt"
	"github.com/bachelor/pkg/config"
	"github.com/bachelor/pkg/model"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestDb(t *testing.T) {
	var (
		client *Db
		err    error
	)
	if err = config.InitConfig(filepath.Join("..", "..")); err != nil {
		t.Error(err)
	}
	if client, err = New(
		viper.GetString("postgres.login"),
		GetEnv("POSTGRES_PASSWORD", "root"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.name"),
	); err != nil {
		t.Error(err)
	}

	filters := []*model.Filter{
		{Filter: "filter1"},
		{Filter: "filter2"},
		{Filter: "filter3"},
	}
	rules := []*model.Rule{
		{Rule: "rule1"},
		{Rule: "rule2"},
		{Rule: "rule3"},
	}
	filtrationRules := []*model.FiltrationRule{
		{
			Filter:         filters[0],
			Rule:           rules[0],
			FilterField:    "field1",
			FilterFunction: "function1",
			FilterValue:    "value1",
		},
		{
			FilterField:    "field2",
			FilterFunction: "function2",
			FilterValue:    "value2",
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

	t.Run("create_rule", func(t2 *testing.T) {
		err := client.CreateRule(rules[0])
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(fmt.Sprintf("filter created successfully: id = %d", rules[0].Id))

	})

	t.Run("create_filtration_rule", func(t2 *testing.T) {
		err := client.CreateFiltrationRule(filtrationRules[0])
		t2.Log(filtrationRules[0].Filter)
		t2.Log(filtrationRules[0].Rule)
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
