package mysql

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
)

// Command model
type Command struct {
	ID           int64
	Cmd          string
	Username     string
	HostID       int64
	Timestamp    int64
	ExitCode     int64
	StdoutLineNo uint64
	StderrLineNo uint64

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

// AddStdoutLine adds line to command_stdout
func (c *Command) AddStdoutLine(ctx context.Context, b []byte) error {
	statement := `
	INSERT INTO command_stdout (id, line_no, line)
	VALUES (?, ?, ?);
	`
	c.StdoutLineNo++
	result, err := c.db.ExecContext(ctx, statement, c.ID, c.StdoutLineNo, b)
	if err != nil {
		return err
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRows != 1 {
		return errors.New("row not inserted")
	}
	log.Debugf("wrote %d stdout lines", c.StdoutLineNo)
	return nil
}

// AddStderrLine adds line to command_stderr
func (c *Command) AddStderrLine(ctx context.Context, b []byte) error {
	statement := `
	INSERT INTO command_stderr (id, line_no, line)
	VALUES (?, ?, ?)
	`
	c.StderrLineNo++
	result, err := c.db.ExecContext(ctx, statement, c.ID, c.StderrLineNo, b)
	if err != nil {
		return err
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRows != 1 {
		return errors.New("row not inserted")
	}
	log.Debugf("wrote %d stderr lines", c.StderrLineNo)
	return nil
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
