package config

// App contains the configuration for the app
type App struct {
	Address string
	Port    string
}

// NewApp returns a new App configuration instance
func NewApp() *App {
	return &App{}
}

// AppGlobal is the global configuration instance for app
var AppGlobal = NewApp()

// AppInit initializes App configuration
func AppInit() {
	Viper.SetEnvPrefix("METAL_APP")
	Viper.SetDefault("address", "0.0.0.0")
	Viper.SetDefault("port", "8000")
	Viper.AutomaticEnv()
	AppGlobal.Address = Viper.GetString("address")
	AppGlobal.Port = Viper.GetString("port")
}
