package config;

import (
	"github.com/spf13/viper"
)

type ConfigKey string;

const (
	ServerPortKey ConfigKey = "server.port"
	IsBootstrapKey = "bootstrap.bootstrap"
	BootstrapServerAddressKey = "bootstrap.bootstrap-address"
)

func setupDefaults() {
	setupDefaultWithKey(ServerPortKey, 80);
	setupDefaultWithKey(IsBootstrapKey, false);
	setupDefaultWithKey(BootstrapServerAddressKey, "127.0.0.1:80");
}

func setupDefaultWithKey(key ConfigKey, value any) {
	viper.SetDefault(string(key), value);
}

func InitializeConfig() error {
	setupDefaults();
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/dope-node")
	return viper.ReadInConfig();
}

func GetString(key ConfigKey) string {
	return viper.GetString(string(key))
}

func GetInt(key ConfigKey) int {
	return viper.GetInt(string(key))
}

func GetBool(key ConfigKey) bool {
	return viper.GetBool(string(key))
}
