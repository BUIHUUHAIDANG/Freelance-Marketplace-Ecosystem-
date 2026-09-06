package config

import(
	"os"
)

type Config struct{
	Port string
}

func LoadConfig() (*Config,error){
	config := &Config{
		Port: os.Getenv("PORT"),
	}
	return config, nil
}