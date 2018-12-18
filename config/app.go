package config

import "sync"

// App contains the configuration for the app
type App struct {
	address string
	port    string
	m       *sync.RWMutex
}

// NewApp returns a new App configuration instance
func NewApp() *App {
	return &App{
		m: &sync.RWMutex{},
	}
}

// GetPort returns the port
func (a *App) GetPort() string {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.port
}

// SetPort retrieves the port
func (a *App) SetPort(port string) {
	a.m.Lock()
	defer a.m.Unlock()
	a.port = port
}

// GetAddress returns the adress
func (a *App) GetAddress() string {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.address
}

// SetAddress retrieves the port
func (a *App) SetAddress(address string) {
	a.m.Lock()
	defer a.m.Unlock()
	a.address = address
}

// AppGlobal is the global configuration instance for app
var AppGlobal = NewApp()

// AppInit initializes App configuration
func AppInit() {
	Viper.SetEnvPrefix("METAL_APP")
	Viper.SetDefault("address", "0.0.0.0")
	Viper.SetDefault("port", "8000")
	Viper.AutomaticEnv()
	AppGlobal.SetAddress(Viper.GetString("address"))
	AppGlobal.SetPort(Viper.GetString("port"))
}
