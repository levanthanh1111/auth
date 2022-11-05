package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

var mapStringLogLevel = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

var consoleLogWriter = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
	w.TimeFormat = time.RFC3339
})

var fileLogWriter = func(fileName string, maxSize, maxAge, maxBackups int) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   fileName,
		MaxBackups: maxBackups, // files
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // days
	}
}

var mapStringLogOutput = map[string]io.Writer{
	"":        consoleLogWriter,
	"console": consoleLogWriter,
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
}

type Config struct {
	Level      string `config:"level"`
	Output     string `config:"output"`
	MaxBackups int    `config:"max_backups"`
	MaxSize    int    `config:"max_size"`
	MaxAge     int    `config:"max_age"`
}

func New(c *Config) zerolog.Logger {
	if c == nil {
		c = &Config{}
	}
	output := mapStringLogOutput[c.Output]
	if output == nil {
		output = fileLogWriter(c.Output, c.MaxSize, c.MaxAge, c.MaxBackups)
	}
	level := mapStringLogLevel[c.Level]
	if level == 0 {
		level = zerolog.DebugLevel
	}
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Logger().
		Level(mapStringLogLevel[c.Level])

	return Logger
}
