package logging

type LoggingConfig struct {
	Level  int    `yaml:"level" env-default:"8"`
	Format string `yaml:"format" env-default:"json"`
	Output string `yaml:"output" env-default:"stdout"`
}
