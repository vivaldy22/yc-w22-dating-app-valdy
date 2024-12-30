package config

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Configuration struct {
	HttpPort       int
	DatabaseMaster Database
	DatabaseSlave  Database
	Redis          redis.Options
	HashSecret     string
	JWTSecret      string
	FeatureFlag    struct {
		EnableDatabaseAutoMigrate bool
	}
	SwipeDailyLimit int64
}

type Database struct {
	Host                  string
	Port                  int
	Name                  string
	Schema                string
	Username              string
	Password              string
	MinIdleConnections    int
	MaxOpenConnections    int
	ConnMaxLifetime       time.Duration
	MaxIdleTimeConnection time.Duration
	DebugMode             bool
}

func LoadConfig() Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := Configuration{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c Configuration) GetHttpPort() string {
	return fmt.Sprintf(":%d", c.HttpPort)
}
