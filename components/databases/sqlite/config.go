package sqlitecomp

type Config struct {
	FilePath    string `yaml:"file_path" env-required:"true"`
	CacheSize   int    `yaml:"cache_size" default:"2000"`
	ForeignKeys bool   `yaml:"foreign_keys" default:"true"`
}
