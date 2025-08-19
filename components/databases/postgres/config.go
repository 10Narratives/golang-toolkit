package pgapp

type Config struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     uint16 `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true"`
}
