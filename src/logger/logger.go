package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

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
	msg := fmt.Sprintf("%s [%s] -- %s\n", timestamp, colorFormater(strings.ToUpper(entry.Level.String())), entry.Message)
	return []byte(msg), nil
}

func Init(level string) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Fatal(err)
	}

	logger.SetFormatter(new(costumeLogFormatter))
	logger.SetLevel(lv)
	logger.Info(fmt.Sprintf("log level: %s", logger.GetLevel().String()))
}

func Get() *logrus.Logger {
	return logger
}
