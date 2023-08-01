package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"smartkid/services/common/context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	// OutputStdout check config(logger.output).
	OutputStdout = "stdout"
	// OutputFile check config(logger.output).
	OutputFile = "file"

	// FormatText check config(logger.format).
	FormatText = "text"
	// FormatJSON check config(logger.format).
	FormatJson = "json"
)

// Config of logger
type Config struct {
	// output can select 'stdout' value or 'file' value. you should must set 'file' key when file value output file.
	Output string `config:"output"`
	// the file path, this is requried if the output is file.
	File  string `config:"file"`
	Level string `config:"level"`
	// format can be text or json
	Format string `config:"format"`
	// disable color default is false, if set to true, the logging show the log with terminal color code
	DisableColors bool `config:"disable_colors"`
	// if expecting line is setted, the code line will be attaching on log whenever called logger
	ExpectingLine bool `config:"expecting_line"`
	// if true we will show log from the payload to console
	RequestPayload bool `config:"request_payload"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Output, validation.Required, validation.In(OutputStdout, OutputStdout)),
		validation.Field(&c.Level, validation.Required, validation.In("panic", "fatal", "error", "warn", "warning", "info", "debug", "trace")),
		validation.Field(&c.Format, validation.Required, validation.In(FormatText, FormatJson)),
		validation.Field(&c.File, validation.Required.When(c.Output == OutputFile)),
	)
}

var Default Logger

func SetDefault(logger Logger) {
	if logger == nil {
		return
	}
	Default = logger
}

// Logger interface
type Logger interface {
	Close() error
	Panic(...interface{})
	Fatal(...interface{})
	Error(...interface{})
	ErrorCtx(context.Context, error, ...interface{})
	// deprecated
	ErrorWrapCtx(context.Context, error, ...interface{}) error
	Warning(...interface{})
	WarningCtx(context.Context, ...interface{})
	Debug(...interface{})
	DebugCtx(context.Context, ...interface{})
	Trace(...interface{})
	TraceCtx(context.Context, ...interface{})
	Info(...interface{})
	InfoCtx(context.Context, ...interface{})
}

type logger struct {
	log           *logrus.Logger
	logfile       *os.File
	expectingLine bool
}

// Close close logfile
func (l *logger) Close() error {
	if l.logfile == nil {
		return nil
	}
	err := l.logfile.Close()
	if err != nil {
		return err
	}
	l.logfile = nil
	return nil
}

func (l *logger) Trace(i ...interface{}) {
	l.withDefaultFields().Trace(i...)
}

func (l *logger) TraceCtx(ctx context.Context, i ...interface{}) {
	l.withDefaultFields(ctx).Trace(i...)
}

func (l *logger) Panic(i ...interface{}) {
	l.withDefaultFields().Panic(i...)
}

func (l *logger) Fatal(i ...interface{}) {
	l.withDefaultFields().Fatal(i...)
}

func (l *logger) Error(i ...interface{}) {
	l.withDefaultFields().Error(i...)
}

func (l *logger) Warning(i ...interface{}) {
	l.withDefaultFields().Warn(i...)
}

func (l *logger) WarningCtx(ctx context.Context, i ...interface{}) {
	l.withDefaultFields(ctx).Warn(i...)
}

func (l *logger) Debug(i ...interface{}) {
	l.withDefaultFields().Debug(i...)
}

func (l *logger) DebugCtx(ctx context.Context, i ...interface{}) {
	l.withDefaultFields(ctx).Debug(i...)
}

func (l *logger) Info(i ...interface{}) {
	l.withDefaultFields().Info(i...)
}

func (l *logger) InfoCtx(ctx context.Context, i ...interface{}) {
	l.withDefaultFields(ctx).Info(i...)
}

func (l *logger) ErrorCtx(ctx context.Context, err error, i ...interface{}) {
	if err != nil {
		l.withDefaultFields(ctx).Error(append(i, " : ", err)...)
	} else {
		l.withDefaultFields(ctx).Error(i...)
	}
}

func (l *logger) ErrorWrapCtx(ctx context.Context, err error, i ...interface{}) error {
	l.withDefaultFields(ctx).Error(append(i, err)...)
	return err
}

func (l *logger) withDefaultFields(ctx ...context.Context) *logrus.Entry {
	locName, locPath := l.getLocationField()
	var fields logrus.Fields
	if len(ctx) > 0 {
		fields = logrus.Fields{locName: locPath}
		// traceId := ctx[0].GetTracerID()
		// if traceId != "" {
		// 	fields["trace"] = traceId
		// }
		userId := ctx[0].GetUserId()
		if userId != "" {
			fields["user_id"] = userId
		}
	} else {
		fields = logrus.Fields{locName: locPath}
	}
	return l.log.WithFields(fields)
}

func (l *logger) getLocationField() (string, string) {
	pc, _, line, _ := runtime.Caller(3)
	frs := runtime.CallersFrames([]uintptr{pc})
	fr, _ := frs.Next()
	var loc string
	if l.expectingLine {
		loc = fmt.Sprintf("%s:%d", fr.Function, line)
	} else {
		loc = fr.Function
	}
	return "location", loc
}

// NewLogger returns new Logger.
// repository: https://github.com/sirupsen/logrus
func NewLogger(c *Config) Logger {
	if c == nil {
		return nil
	}
	var err error
	var file *os.File

	// new logrus.
	log := logrus.New()

	// set output.
	switch c.Output {
	case OutputStdout: // output: stdout
		log.Out = os.Stdout
		logrus.SetOutput(os.Stdout)
	case OutputFile: // output: file
		file, err = os.OpenFile(c.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			_, fileName := path.Split(c.File)
			file, err = os.OpenFile(fmt.Sprintf("%s/%s", os.Getenv("API_DIR"), fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				return nil
			}
		}
		log.Out = file
		logrus.SetOutput(file)
	default:
		log.Out = os.Stdout
		logrus.SetOutput(os.Stdout)
	}

	// set formatter.
	switch c.Format {
	case FormatText:
		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		// Set specific colors for prefix and timestamp
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle:    "blue+b",
			TimestampStyle: "grey+h",
		})
		if c.Output != OutputStdout || c.DisableColors {
			formatter.DisableColors = true
		}
		log.SetFormatter(formatter)
		logrus.SetFormatter(formatter)
	case FormatJson:
		formatter := &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "func",
			},
		}
		logrus.SetFormatter(formatter)
		log.Formatter = formatter
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
		log.Formatter = &logrus.TextFormatter{}
	}

	level := c.Level
	// set level.
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return nil
	}
	if level == "" {
		level = "info"
	}

	log.Info("log level: ", lv)
	log.SetLevel(lv)
	logrus.SetLevel(lv)

	return &logger{log: log, logfile: file, expectingLine: c.ExpectingLine}
}

func Factory(injector *do.Injector) (Logger, error) {
	return NewLogger(do.MustInvoke[*Config](injector)), nil
}
