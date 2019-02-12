package mysql

import (
	"context"
	"database/sql"
)

// Command model
type Command struct {
	ID        int64
	Cmd       string
	Username  string
	HostID    int64
	Timestamp int64
	ExitCode  int64
}

// CommandOpt is an option for a NewCommand
type CommandOpt func(*Command)

// WithCommandID adds an id to the NewCommand
func WithCommandID(id int64) CommandOpt {
	return func(c *Command) {
		c.ID = id
	}
}

// WithExitCode adds an exit code to the NewCommand
func WithExitCode(exitCode int64) CommandOpt {
	return func(c *Command) {
		c.ExitCode = exitCode
	}
}

// NewCommand returns a new Command
func NewCommand(db *sql.DB, cmd string, username, hostID string, timestamp int64, opts ...CommandOPt) *Command {
	c := &Command{
		Cmd:       cmd,
		Timestamp: timestamp,
		db:        db,
	}
	for _, fn := range opts {
		fn(c)
	}
	return c
}

// Create a Command record
func (c *Command) Create(ctx context.Context) error {
	query := `
  INSERT INTO command (cmd, username, host_id, timestamp)
  VALUES (?, ?, ?, ?, ?, ?);
  `
	result, err := db.ExecContext(ctx, query, c.Cmd, c.Username, c.HostID, c.Timestamp)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	c.ID = id
	return err
}
