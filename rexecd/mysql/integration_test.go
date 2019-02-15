// +build mysql

package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql/migration"
	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"github.com/metal-go/metal/test/lib"
)

func TestHost(t *testing.T) {
	config.RexecdInit()
	if err := migration.New().Run(); err != nil {
		t.Error(err)
	}
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("create", func(t *testing.T) {
		host := NewHost(db)

		id, err := host.Create(context.Background(), "test-host", WithHostKeyType(proto_rexecd.KeyType_rsa_sha2_512),
			WithHostPublicKey([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		if id != 1 {
			t.Errorf("expected lastInsertID to be 1, got %d", id)
		}
	})
	t.Run("Read", func(t *testing.T) {
		host := NewHost(db)
		if err := host.Read(context.Background(), "test-host"); err != nil {
			t.Fatal(err)
		}
		if host.FQDN != "test-host" {
			t.Errorf("expected test-host, got %s", host.FQDN)
		}
	})
}

func TestUser(t *testing.T) {
	config.RexecdInit()
	if err := migration.New().Run(); err != nil {
		t.Error(err)
	}
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("create", func(t *testing.T) {
		user := NewUser(db)
		if err := user.Create(context.Background(), "test-user", WithUserFirstName("test"), WithUserLastName("user")); err != nil {
			t.Fatal(err)
		}
		var size int
		row := db.QueryRow("SELECT COUNT(*) FROM user;")
		if err := row.Scan(&size); err != nil {
			t.Fatal(err)
		}
		if size != 1 {
			t.Fatal(size)
		}
	})
	t.Run("Read", func(t *testing.T) {
		user := NewUser(db)
		if err := user.Read(context.Background(), "test-user"); err != nil {
			t.Fatal(err)
		}
		// result, err := db.Query("SELECT * FROM user;")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// var username string
		// var firstName string
		// var lastName string
		// var admin bool
		//
		// for result.Next() {
		// 	var username string
		// 	if err := result.Scan(&username, &firstName, &lastName, &admin); err != nil {
		// 		t.Fatal(err)
		// 	}
		// }
		// fmt.Println(username, firstName, lastName, admin)
		if user.FirstName != "test" {
			t.Errorf("expected test, got %s", user.FirstName)
		}
	})
}

func TestMain(t *testing.M) {
	lib.MySQLTestMain(t)
}
