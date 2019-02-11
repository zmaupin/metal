package migration

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql"
)

const initDB = `
CREATE DATABASE IF NOT EXISTS rexecd;
`

const initMigrationsTable = `
CREATE TABLE IF NOT EXISTS migration (
  id SERIAL,
	data BLOB,
  PRIMARY KEY (id)
);
`

var migrations = []string{}

// Migrate handles migrations and initialization
type Migrate struct {
	ctx     context.Context
	cancel  func()
	timeout time.Duration
}

// MigrateOpt is an option for Migrate
type MigrateOpt func(*Migrate)

// WithTimeout adds timeout to Migrate
func WithTimeout(t time.Duration) MigrateOpt {
	return func(m *Migrate) {
		m.timeout = t
	}
}

// New returns a pointer to a new Migrate
func New(opts ...MigrateOpt) *Migrate {
	m := &Migrate{}

	for _, fn := range opts {
		fn(m)
	}

	if m.timeout == 0 {
		m.timeout = time.Second * 60
	}

	m.ctx, m.cancel = context.WithTimeout(context.Background(), m.timeout)
	return m
}

// Run executes the migration
func (m *Migrate) Run() error {
	defer m.cancel()
	db, err := sql.Open("mysql", config.RexecdGlobal.DataSourceName)
	if err != nil {
		return err
	}
	if _, err = db.ExecContext(m.ctx, initDB); err != nil {
		return err
	}
	db.Close()

	db, err = sql.Open("mysql", config.RexecdGlobal.DataSourceName+"rexecd")
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err = db.ExecContext(m.ctx, initMigrationsTable); err != nil {
		return err
	}

	rows, err := db.QueryContext(m.ctx, "SELECT id FROM migration ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return err
		}
	}

	if id == len(migrations) {
		return nil
	}

	for _, mig := range migrations[id:] {
		if err := mysql.ExecuteSQL(m.ctx, db, mig); err != nil {
			return err
		}
		if _, err := db.QueryContext(m.ctx, "INSERT INTO migration (data) VALUES (?);\n", []byte(mig)); err != nil {
			return err
		}
	}
	log.Info("databse initialized")
	return nil
}

func init() {
	migrations = append(migrations, _0)
}
