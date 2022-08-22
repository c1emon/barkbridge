package ilogger

import (
	"fmt"
	"strings"
	"time"

	"github.com/c1emon/barkbridge/utils"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
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
	msg := fmt.Sprintf("%s [%s] -- %s\n",
		timestamp,
		colorFormater(strings.ToUpper(entry.Level.String())),
		entry.Message)
	if entry.Data != nil && len(entry.Data) > 0 {
		msg = fmt.Sprintf("%swith values:\n%s\n", msg, utils.PrettyMarshal(entry.Data))
	}

	return []byte(msg), nil
}

func Init(level string) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetFormatter(new(costumeLogFormatter))
	logrus.SetLevel(lv)
	logrus.Info(fmt.Sprintf("log level: %s", logrus.GetLevel().String()))
}
