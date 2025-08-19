package grpcsrv

type Config struct {
	Address string `yaml:"address" env-required:"true"`
	Port    int    `yaml:"port" env-required:"true"`
}
