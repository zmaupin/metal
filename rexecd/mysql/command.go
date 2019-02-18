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

	db *sql.DB
}

// CommandOpt is an option for a NewCommand
type CommandOpt func(*Command)

// NewCommand returns a new Command
func NewCommand(db *sql.DB) *Command {
	c := &Command{db: db}
	return c
}

// Create a Command record
func (c *Command) Create(ctx context.Context, cmd string, username, fqdn string, timestamp int64, opts ...CommandOpt) error {
	for _, fn := range opts {
		fn(c)
	}
	c.Cmd = cmd
	c.Username = username
	c.Timestamp = timestamp

	query := `
	SELECT id FROM host WHERE fqdn = ?;
	`
	var hostID int64
	row := c.db.QueryRowContext(ctx, query, fqdn)
	if err := row.Scan(&hostID); err != nil {
		return err
	}
	c.HostID = hostID

	statement := `
	INSERT INTO command (cmd, username, host_id, timestamp)
	VALUES (?, ?, ?, ?);
  `
	result, err := c.db.ExecContext(ctx, statement, cmd, username, hostID, timestamp)
	if err != nil {
		return err
	}
	cmdID, err := result.LastInsertId()
	c.ID = cmdID
	return err
}

// AddStdout adds stdout
func (c *Command) AddStdout(ctx context.Context, b []byte) error {
	statement := `
	UPDATE command
	SET stdout = ?
	WHERE id = ?;
	`
	_, err := c.db.ExecContext(ctx, statement, b, c.ID)
	return err
}

// AddStderr adds stderr
func (c *Command) AddStderr(ctx context.Context, b []byte) error {
	statement := `
	UPDATE command
	SET stderr = ?
	WHERE id = ?;
	`
	_, err := c.db.ExecContext(ctx, statement, b, c.ID)
	return err
}

// SetExitCode sets the exit code on command
func (c *Command) SetExitCode(ctx context.Context, exitCode int64) error {
	c.ExitCode = exitCode
	statement := `
	UPDATE command
	SET exit_code = ?
	WHERE id = ?;
	`
	_, err := c.db.ExecContext(ctx, statement, exitCode, c.ID)
	return err
}
