package migration

import (
	"context"
	"database/sql"
	"time"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql"
)

const initDB = `
CREATE DATABASE IF NOT EXISTS rexecd;
`

const initMigrationsTable = `
CREATE TABLE IF NOT EXISTS migrations (
  id int NOT NULL AUTO_INCREMENT,
  succeeded tinyint, 
  PRIMARY KEY (id));
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
	if err = mysql.ExecuteSQL(m.ctx, db, initDB); err != nil {
		return err
	}

	db, err = sql.Open("mysql", config.RexecdGlobal.DataSourceName+"rexecd")
	if err != nil {
		return err
	}

	return mysql.ExecuteSQL(m.ctx, db, initMigrationsTable)
}
