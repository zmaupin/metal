package migration

import (
	"context"
	"database/sql"
	"time"

	"github.com/metal-go/metal/config"
)

const initDB = `
CREATE DATABASE IF NOT EXISTS rexecd;
`

const initMigrationsTable = `
CREATE TABLE IF NOT EXISTS migration (
  id int NOT NULL AUTO_INCREMENT,
  success BOOLEAN,
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
	if _, err = db.ExecContext(m.ctx, initDB); err != nil {
		return err
	}

	db, err = sql.Open("mysql", config.RexecdGlobal.DataSourceName+"rexecd")
	if err != nil {
		return err
	}

	if _, err = db.ExecContext(m.ctx, initMigrationsTable); err != nil {
		return err
	}

	rows, err := db.QueryContext(m.ctx, "SELECT * FROM migration;\n")
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		var success int
		if err := rows.Scan(&id, &success); err != nil {
			return err
		}
		if success == 1 {
			continue
		}
		if _, err := db.ExecContext(m.ctx, migrations[id-1]); err != nil {
			return err
		}
		db.QueryContext(m.ctx, "INSERT INTO migration (id, success) VALUES (?, TRUE);\n", id)
	}

	if id == len(migrations) {
		return nil
	}

	for i, mig := range migrations[id-1:] {
		db.QueryContext(m.ctx, "INSERT INTO migration VALUES (FALSE);\n")
		if _, err := db.ExecContext(m.ctx, mig); err != nil {
			return err
		}
		db.QueryContext(m.ctx, "INSERT INTO migration (id, success) VALUES (?, TRUE);\n", i+1)
	}
	return nil
}
