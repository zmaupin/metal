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
func NewUser(db *sql.DB, username string, opts ...UserOpt) *User {
	u := &User{
		Username: username,
		db:       db,
	}
	for _, fn := range opts {
		fn(u)
	}
	return u
}

// Create a new User
func (u *User) Create(ctx context.Context) error {
	if _, err := u.Read(ctx, u.Username); err != nil {
		return err
	}
	statement := `
	INSERT INTO user (username, first_name, last_name, admin)
	VALUES (?, ?, ?, ?, ?, ?);
  `
	_, err := u.db.ExecContext(ctx, statement, u.Username, u.FirstName, u.LastName, u.Admin)
	return err
}

// Read a User
func (u *User) Read(ctx context.Context, username string) (*User, error) {
	var userName string
	var firstName string
	var lastName string

	query := `
	SELECT (username, first_name, last_name, admin)
	FROM user
	WHERE username = ?;
  `
	row := u.db.QueryRowContext(ctx, query, username)
	err := row.Scan()

	return &User{
		Username:  userName,
		FirstName: firstName,
		LastName:  lastName,
	}, err
}
