// +build integration

package migration

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/test/lib"
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
	t.Run("test_count", func(t *testing.T) {
		rows, err := db.Query("SELECT COUNT(*) FROM migration;\n")
		if err != nil {
			t.Fatal(err)
		}
		var count int
		for rows.Next() {
			if err = rows.Scan(&count); err != nil {
				t.Fatal(err)
			}
			if count != 1 {
				t.Errorf("expected count to be 1, got %d", count)
			}
		}
	})
	t.Run("test_migration_data", func(t *testing.T) {
		rows, err := db.Query("SELECT id, data FROM migration;\n")
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var id int
			var data []byte

			if err := rows.Scan(&id, &data); err != nil {
				t.Fatal(err)
			}

			if string(data) != _0 {
				t.Errorf("expected %s, got %s", _0, string(data))
			}
		}
	})
	t.Run("test _01", func(t *testing.T) {
		_, err := db.Query("SELECT * FROM host;\n")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMain(t *testing.M) {
	lib.MySQLTestMain(t)
}
