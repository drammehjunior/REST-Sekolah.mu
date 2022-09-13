package config

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
}

func LoadConfig() (Config, error) {
	return Config{
		DBHost:     "localhost",
		DBName:     "test",
		DBUser:     "root",
		DBPort:     "3306",
		DBPassword: "",
	}, nil
}
