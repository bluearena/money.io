package config

import (
    "github.com/kelseyhightower/envconfig"
    "github.com/gromnsk/money.io/pkg/logger"
)

const (
    // SERVICENAME contains a service name prefix which used in ENV variables
    SERVICENAME = "MONEYIO"
)

// Config contains ENV variables
type Config struct {
    // Local service host
    LocalHost string `split_words:"true"`
    // Local service port
    LocalPort int `split_words:"true"`
    // Logging level in logger.Level notation
    LogLevel logger.Level `split_words:"true"`
    // Telegram bot API token
    TelegramToken string `split_words:"true"`
}

// DBConfig set configuration of database
type DBConfig struct {
    Host     string `required:"true"`
    Port     string `required:"true"`
    Database string `required:"true"`
    Username string `required:"true"`
    Password string `required:"true"`
}

// Load settles ENV variables into Config structure
func (c *Config) Load(serviceName string) error {
    return envconfig.Process(serviceName, c)
}
