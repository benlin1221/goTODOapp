package config

type DatabaseDriver struct {
	Driver      string `yaml:"driver" env:"DB_DRIVER"`
	Host        string `yaml:"host" env:"DB_HOST"`
	Username    string `yaml:"username" env:"DB_USER"`
	Password    string `yaml:"password" env:"DB_PASS"`
	DBName      string `yaml:"db_database" env:"DB_DATABASE"`
	Port        int    `yaml:"port" env:"DB_PORT"`
	Connections int    `yaml:"connections" env:"DB_CONNECTIONS"`
}

type DatabaseConfig struct {
	Drivers map[string]DatabaseDriver `yaml:"drivers"`
	Default DatabaseDriver            `yaml:"default" env:"DEFAULT_DB_DRIVER"`
}
