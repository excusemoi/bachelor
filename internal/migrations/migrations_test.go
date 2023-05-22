package migrations

import (
	"testing"
)

func TestMigrations(t *testing.T) {

	if err := Run(); err != nil {
		t.Fatal(err)
	}

}
