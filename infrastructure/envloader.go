package infrastructure

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gobot/entities"
)

func LoadConfigFromEnv() (entities.Config, error) {
	err := godotenv.Load(".env")
	var config entities.Config
	err = envconfig.Process("", &config)
	if err != nil {
		return entities.Config{}, err
	}
	return config, nil
}
