package migrations

import (
	"testing"
)

func TestMigrations(t *testing.T) {

	if err := Run(); err != nil {
		t.Fatal(err)
	}

	//t.Run("get_filtration_rule", func(t2 *testing.T) {
	//	fr, err := client.GetByID(1)
	//	if err != nil {
	//		t2.Error(err)
	//	}
	//	t2.Log(fr)
	//})
	//fmt.Println(client)
}
