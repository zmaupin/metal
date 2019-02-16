package config

// Rexecd represents the global configuration for Rexecd
type Rexecd struct {
	Cluster           []string
	ServerType        string
	Address           string
	Port              string
	DataSourceName    string
	CommandTimeoutSec int
}

// NewRexecd returns a new Rexecd configuration instance
func NewRexecd() *Rexecd {
	return &Rexecd{}
}

// RexecdGlobal is the global configuration for Rexecd
var RexecdGlobal = NewRexecd()

// RexecdInit initializes the global configuration for rexecd.
func RexecdInit() {
	Viper.SetEnvPrefix("METAL_REXECD")
	Viper.SetDefault("cluster", "")
	Viper.SetDefault("server_type", "mysql")
	Viper.SetDefault("address", "0.0.0.0")
	Viper.SetDefault("port", "9000")
	Viper.SetDefault("command_timeout_sec", 300)
	Viper.AutomaticEnv()
	RexecdGlobal.Cluster = arrayValue(Viper.GetString("cluster"))
	RexecdGlobal.ServerType = Viper.GetString("server_type")
	RexecdGlobal.Address = Viper.GetString("address")
	RexecdGlobal.Port = Viper.GetString("port")
	RexecdGlobal.DataSourceName = Viper.GetString("data_source_name")
	RexecdGlobal.CommandTimeoutSec = Viper.GetInt("command_timeout_sec")
}
