// +build mysql

package migration

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/metal-go/metal/config"
)

func TestMigrate(t *testing.T) {
	config.RexecdInit()
	if err := New().Run(); err != nil {
		t.Error(err)
	}
}
