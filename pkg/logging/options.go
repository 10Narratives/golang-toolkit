package logging

type options struct {
	level  int
	format string
	output string
}

func defaultOptions() *options {
	return &options{
		level:  8,
		format: "json",
		output: "stdout",
	}
}

type LoggerOption func(*options)

func WithLevel(level int) LoggerOption {
	return func(lo *options) {
		lo.level = level
	}
}

func WithFormat(format string) LoggerOption {
	return func(lo *options) {
		lo.format = format
	}
}

func WithOutput(output string) LoggerOption {
	return func(lo *options) {
		lo.output = output
	}
}
