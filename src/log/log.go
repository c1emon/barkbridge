package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type Level logrus.Level

const (
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

var cstZone = time.FixedZone("GMT", 8*3600)

// CostumeLogFormatter Custom log format definition
type costumeLogFormatter struct{}

// Format log format
func (s *costumeLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	var colorFormater func(a ...interface{}) string
	switch entry.Level {
	case logrus.DebugLevel:
		colorFormater = color.New(color.FgHiYellow).SprintFunc()
	case logrus.InfoLevel:
		colorFormater = color.New(color.FgGreen).SprintFunc()
	case logrus.WarnLevel:
		colorFormater = color.New(color.FgYellow).SprintFunc()
	default:
		colorFormater = color.New(color.FgRed).SprintFunc()
	}

	timestamp := time.Now().In(cstZone).Format("2006-01-02 15:04:05.999")
	msg := fmt.Sprintf("%s [%s] --- %s\n", timestamp, colorFormater(strings.ToUpper(entry.Level.String())), entry.Message)
	return []byte(msg), nil
}

func Init(level Level) {
	logger.SetFormatter(new(costumeLogFormatter))
	logger.SetLevel(logrus.Level(level))
	Info(fmt.Sprintf("log level: %s", logger.GetLevel().String()))
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}
