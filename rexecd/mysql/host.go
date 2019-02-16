package mysql

import (
	"context"
	"database/sql"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"github.com/metal-go/metal/rexecd"
)

// Host model
type Host struct {
	ID        int64
	FQDN      string
	Port      string
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

// WithHostPort adds the Port to the Host
func WithHostPort(port string) HostOpt {
	return func(h *Host) {
		h.Port = port
	}
}

// WithHostPublicKey adds the public key to the Host
func WithHostPublicKey(key []byte) HostOpt {
	return func(h *Host) {
		h.PublicKey = key
	}
}

// WithHostKeyType adds the keytype to the Host
func WithHostKeyType(keyType proto_rexecd.KeyType) HostOpt {
	return func(h *Host) {
		h.KeyType = keyType
	}
}

// NewHost returns a new host
func NewHost(db *sql.DB) *Host {
	h := &Host{db: db}
	if h.Port == "" {
		h.Port = "22"
	}
	return h
}

// Create a new Host
func (h *Host) Create(ctx context.Context, fqdn string, opts ...HostOpt) (id int64, err error) {
	h.FQDN = fqdn
	for _, fn := range opts {
		fn(h)
	}
	statement := `
	INSERT INTO host (fqdn, port, public_key, key_type) VALUES (?, ?, ?, ?);
  `
	result, err := h.db.ExecContext(ctx, statement, h.FQDN, h.Port, h.PublicKey, rexecd.KeyTypeKey(h.KeyType))
	if err != nil {
		return int64(0), err
	}
	return result.LastInsertId()
}

// Read gets Host info
func (h *Host) Read(ctx context.Context, fqdn string) error {
	query := `
	SELECT id, fqdn, port, public_key, key_type FROM host WHERE fqdn = ?;
  `
	row := h.db.QueryRowContext(ctx, query, fqdn)

	var id int64
	var port string
	var publicKey []byte
	var keyType string

	if err := row.Scan(&id, &fqdn, &port, &publicKey, &keyType); err != nil {
		return err
	}
	h.ID = id
	h.Port = port
	h.PublicKey = publicKey
	h.KeyType = proto_rexecd.KeyType(proto_rexecd.KeyType_value[keyType])
	h.FQDN = fqdn
	return nil
}
