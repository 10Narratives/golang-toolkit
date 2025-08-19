package logging

type LoggingConfig struct {
	Level  string `yaml:"level" env-default:"8"`
	Format string `yaml:"format" env-default:"json"`
	Output string `yaml:"output" env-default:"stdout"`
}
