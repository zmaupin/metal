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
		host := NewHost(
			db, WithHostFQDN("test-host"), WithHostKeyType(proto_rexecd.KeyType_rsa_sha2_512),
			WithHostPublicKey([]byte{}))

		if _, err := host.Create(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMain(t *testing.M) {
	lib.MySQLTestMain(t)
}
