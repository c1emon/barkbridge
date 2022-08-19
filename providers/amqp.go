package providers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func Dail() {
	amqp.Dial("")
	logrus.Debug()
}
