package migration

import (
	"testing"

	"github.com/xwb1989/sqlparser"
)

func TestInitDB(t *testing.T) {
	_, err := sqlparser.Parse(initDB)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestInitMigrationsTable(t *testing.T) {
	_, err := sqlparser.Parse(initMigrationsTable)
	if err != nil {
		t.Error(err.Error())
	}
}
