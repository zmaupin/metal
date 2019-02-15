package mysql

import (
	"context"
	"database/sql"
)

// User model
type User struct {
	Username  string
	FirstName string
	LastName  string
	Admin     bool
	db        *sql.DB
}

// UserOpt option for a new User
type UserOpt func(*User)

// WithUsername sets the optional FirstName
func WithUsername(name string) UserOpt {
	return func(u *User) {
		u.Username = name
	}
}

// WithUserFirstName sets the optional FirstName
func WithUserFirstName(name string) UserOpt {
	return func(u *User) {
		u.FirstName = name
	}
}

// WithUserLastName sets the optional LastName
func WithUserLastName(name string) UserOpt {
	return func(u *User) {
		u.LastName = name
	}
}

// WithUserAdmin sets the optional Admin member to true
func WithUserAdmin(b bool) UserOpt {
	return func(u *User) {
		u.Admin = b
	}
}

// NewUser returns a new User
func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

// Create a new User
func (u *User) Create(ctx context.Context, username string, opts ...UserOpt) error {
	WithUsername(username)(u)
	for _, fn := range opts {
		fn(u)
	}

	statement := `
	INSERT INTO user (username, first_name, last_name, admin)
	VALUES (?, ?, ?, ?);
  `
	_, err := u.db.ExecContext(ctx, statement, u.Username, u.FirstName, u.LastName, u.Admin)
	return err
}

// Read a User
func (u *User) Read(ctx context.Context, username string) error {
	var firstName string
	var lastName string
	var admin bool

	query := `
	SELECT first_name, last_name, admin
	FROM user
	WHERE username = ?;
  `
	row := u.db.QueryRowContext(ctx, query, username)
	err := row.Scan(&firstName, &lastName, &admin)

	u.FirstName = firstName
	u.LastName = lastName
	u.Admin = admin
	return err
}
