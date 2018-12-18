package config

import "strings"
import "sync"

// Rexecd represents the global configuration for Rexecd
type Rexecd struct {
	cluster    []string
	serverType string
	address    string
	port       string
	m          *sync.RWMutex
}

// NewRexecd returns a new Rexecd configuration instance
func NewRexecd() *Rexecd {
	return &Rexecd{
		m: &sync.RWMutex{},
	}
}

// SetCluster in Rexecd
func (r *Rexecd) SetCluster(c []string) {
	r.m.Lock()
	defer r.m.Unlock()
	r.cluster = c
}

// GetCluster returns the cluster information
func (r *Rexecd) GetCluster() []string {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.cluster
}

// GetServerType retrieves the serverType
func (r *Rexecd) GetServerType() string {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.serverType
}

// SetServerType sets the serverType
func (r *Rexecd) SetServerType(serverType string) {
	r.m.Lock()
	defer r.m.Unlock()
	r.serverType = strings.ToLower(serverType)
}

// GetPort returns the port
func (r *Rexecd) GetPort() string {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.port
}

// SetPort retrieves the port
func (r *Rexecd) SetPort(port string) {
	r.m.Lock()
	defer r.m.Unlock()
	r.port = port
}

// GetAddress returns the adress
func (r *Rexecd) GetAddress() string {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.address
}

// SetAddress retrieves the port
func (r *Rexecd) SetAddress(address string) {
	r.m.Lock()
	defer r.m.Unlock()
	r.address = address
}

// RexecdGlobal is the global configuration for Rexecd
var RexecdGlobal = NewRexecd()

// RexecdInit initializes the global configuration for rexecd.
func RexecdInit() {
	Viper.SetEnvPrefix("METAL_REXECD")
	Viper.SetDefault("cluster", "")
	Viper.SetDefault("server_type", "memory")
	Viper.SetDefault("address", "0.0.0.0")
	Viper.SetDefault("port", "9000")
	Viper.AutomaticEnv()
	RexecdGlobal.SetCluster(arrayValue(Viper.GetString("cluster")))
	RexecdGlobal.SetServerType(Viper.GetString("server_type"))
	RexecdGlobal.SetAddress(Viper.GetString("address"))
	RexecdGlobal.SetPort(Viper.GetString("port"))
}
