package config

type Config struct {
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"vechnonetot"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"dimanik98"`
	DBName     string `envconfig:"DB_NAME" default:"servicedb"`
}

func NewConfig() *Config {
	return &Config{}
}
