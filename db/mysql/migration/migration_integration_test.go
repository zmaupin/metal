// +build mysql

package migration

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/metal-go/metal/config"
)

func TestMigrate(t *testing.T) {
	config.RexecdInit()
	if err := New().Run(); err != nil {
		t.Error(err)
	}
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query("SELECT COUNT(id) FROM migration;\n")
	defer rows.Close()
	if err != nil {
		t.Fatal(err)
	}
	var count int
	rows.Next()
	if err = rows.Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected count to be 1, got %d", count)
	}
}
