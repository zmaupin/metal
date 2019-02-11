package mysql

import (
	"context"
	"database/sql"
	"fmt"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// Host model
type Host struct {
	ID         int64
	FQDN       string
	Port       int64
	PrivateKey []byte
	PublicKey  []byte
	KeyType    proto_rexecd.KeyType
	db         *sql.DB
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
func NewHost(db *sql.DB, fqdn string, port int64, privateKey []byte, publicKey []byte, keyType proto_rexecd.KeyType, opts ...HostOpt) *Host {
	h := &Host{
		FQDN:       fqdn,
		Port:       port,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		KeyType:    keyType,
		db:         db,
	}
	for _, fn := range opts {
		fn(h)
	}
	return h
}

// Create a new Host
func (h *Host) Create(ctx context.Context) (id int64, err error) {
	statement := `
  INSERT INTO host (fqdn, port, private_key, public_key) VALUES (?, ?, ?, ?);
  `
	result, err := h.db.ExecContext(
		ctx, statement, h.FQDN, h.Port, h.PrivateKey, h.PublicKey,
		proto_rexecd.KeyType_name[int32(h.KeyType)])

	if err != nil {
		return int64(0), err
	}
	return result.LastInsertId()
}

// Read gets Host info
func (h *Host) Read(ctx context.Context, id int64) (*Host, error) {
	query := `
  SELECT (fqdn, port, private_key, public_key) FROM host WHERE id = ?;
  `
	rows, err := h.db.QueryContext(ctx, query, id)
	if err != nil {
		return &Host{}, err
	}
	for rows.Next() {
		var fqdn string
		var port int64
		var privateKey []byte
		var publicKey []byte

		err = rows.Scan(&fqdn, &port, &privateKey, &publicKey)
		if err != nil {
			return &Host{}, err
		}
		return &Host{ID: id, FQDN: fqdn, PrivateKey: privateKey, PublicKey: publicKey}, nil
	}
	return &Host{}, fmt.Errorf("host not found for id %d", id)
}
