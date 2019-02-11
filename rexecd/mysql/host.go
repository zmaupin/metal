package mysql

import (
	"context"
	"database/sql"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// Host model
type Host struct {
	ID        int64
	FQDN      string
	Port      int64
	PublicKey []byte
	KeyType   proto_rexecd.KeyType
	db        *sql.DB
}

// HostOpt is an option for a new Host
type HostOpt func(*Host)

// WithHostID returns a HostOpt with a configured ID
func WithHostID(id int64) HostOpt {
	return func(h *Host) {
		h.ID = id
	}
}

// NewHost returns a new host
func NewHost(db *sql.DB, fqdn string, port int64, publicKey []byte, keyType proto_rexecd.KeyType, opts ...HostOpt) *Host {
	h := &Host{
		FQDN:      fqdn,
		Port:      port,
		PublicKey: publicKey,
		KeyType:   keyType,
		db:        db,
	}
	for _, fn := range opts {
		fn(h)
	}
	return h
}

// Create a new Host
func (h *Host) Create(ctx context.Context) (id int64, err error) {
	statement := `
  INSERT INTO host (fqdn, port, public_key, key_type) VALUES (?, ?, ?, ?);
  `
	result, err := h.db.ExecContext(
		ctx, statement, h.FQDN, h.Port, h.PublicKey,
		proto_rexecd.KeyType_name[int32(h.KeyType)])

	if err != nil {
		return int64(0), err
	}
	return result.LastInsertId()
}

// Read gets Host info
func (h *Host) Read(ctx context.Context, id int64) (*Host, error) {
	query := `
  SELECT (fqdn, port, public_key) FROM host WHERE id = ?;
  `
	row := h.db.QueryRowContext(ctx, query, id)

	var fqdn string
	var port int64
	var publicKey []byte

	err := row.Scan(&fqdn, &port, &publicKey)
	return &Host{ID: id, FQDN: fqdn, PublicKey: publicKey}, err
}
