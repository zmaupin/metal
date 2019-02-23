package config

// Rexecd represents the global configuration for Rexecd
type Rexecd struct {
	Address           string
	APIAddress        string
	CommandTimeoutSec int
	Cluster           []string
	DataSourceName    string
	Port              string
	ServerType        string
	KafkaAddress      []string
	KafkaVersion      string
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
	Viper.SetDefault("api_address", ":8080")
	Viper.SetDefault("port", "9000")
	Viper.SetDefault("command_timeout_sec", 300)
	Viper.SetDefault("kafka_version", "2.12-2.1.0")
	Viper.SetDefault("kafka_address", "")
	Viper.AutomaticEnv()
	RexecdGlobal.Address = Viper.GetString("address")
	RexecdGlobal.APIAddress = Viper.GetString("api_address")
	RexecdGlobal.Cluster = arrayValue(Viper.GetString("cluster"))
	RexecdGlobal.CommandTimeoutSec = Viper.GetInt("command_timeout_sec")
	RexecdGlobal.DataSourceName = Viper.GetString("data_source_name")
	RexecdGlobal.KafkaAddress = arrayValue(Viper.GetString("kafka_address"))
	RexecdGlobal.KafkaVersion = Viper.GetString("kafka_version")
	RexecdGlobal.Port = Viper.GetString("port")
	RexecdGlobal.ServerType = Viper.GetString("server_type")
}
