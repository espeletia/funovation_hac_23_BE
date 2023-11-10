package config

import "github.com/spf13/viper"

type Config struct {
	Server           ServerConfig
	Service          ServiceConfig
	DBConfig         DBConfig
	MigrationsConfig MigrationsConfig
	S3Config         S3Config
	EncodingConfig   EncodingConfig
}

func LoadConfig() *Config {
	return &Config{
		Server:           loadServerConfig(),
		Service:          loadServiceConfig(),
		DBConfig:         loadDbConfig(),
		MigrationsConfig: loadMigrationsConfig(),
		S3Config:         loadS3Config(),
		EncodingConfig:   loadEncodingConfig(),
	}
}

func configViper(configName string) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configurations/")
	return v
}
