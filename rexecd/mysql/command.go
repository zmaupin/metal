package mysql

import (
	"context"
	"database/sql"
	"time"
)

// Command model
type Command struct {
	ID        int64
	Cmd       string
	Username  string
	HostID    int64
	Timestamp int64
	ExitCode  int64
	db        *sql.DB
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
func NewCommand(db *sql.DB, cmd string, username string, hostID, timestamp int64, opts ...CommandOpt) *Command {
	c := &Command{
		Cmd:       cmd,
		Username:  username,
		HostID:    hostID,
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
	result, err := c.db.ExecContext(ctx, query, c.Cmd, c.Username, c.HostID, c.Timestamp)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	c.ID = id
	return err
}

// AddStdoutLine adds line to command_stdout
func (c *Command) AddStdoutLine(ctx context.Context, b []byte) error {
	statement := `
	INSERT INTO command_stdout (id, timestamp, line)
	VALUES (?, ?, ?);
	`
	_, err := c.db.ExecContext(ctx, statement, c.ID, time.Now().Unix(), b)
	return err
}

// AddStderrLine adds line to command_stderr
func (c *Command) AddStderrLine(ctx context.Context, b []byte) error {
	statement := `
	INSERT INTO command_stderr (id, timestamp, line)
	VALUES (?, ?, ?)
	`
	_, err := c.db.ExecContext(ctx, statement, c.ID, time.Now().Unix(), b)
	return err
}

// SetExitCode sets the exit code on command
func (c *Command) SetExitCode(ctx context.Context, exitCode int64) error {
	statement := `
	INSERT INTO command (exit_code) VALUES (?) WHERE id = ?;
	`
	_, err := c.db.ExecContext(ctx, statement, exitCode, c.ID)
	return err
}
