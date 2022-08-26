package providers

import (
	"fmt"
	"sync"

	"github.com/c1emon/barkbridge/barkserver"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AmqpProvider struct {
	Addr      string
	ProvideCh chan barkserver.Message
	stopCh    chan any
	wg        *sync.WaitGroup
}

func NewAmqpProvider(addr string) *AmqpProvider {
	p := &AmqpProvider{
		Addr:      addr,
		wg:        &sync.WaitGroup{},
		ProvideCh: make(chan barkserver.Message),
		stopCh:    make(chan any),
	}

	return p
}

func (p *AmqpProvider) Start() {
	conn, err := amqp.Dial(p.Addr)
	if err != nil {
		logrus.WithField("error", err).Panic("failed dial amqp server")
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.WithField("error", err).Fatal("failed open channel")
	}

	exName := "amq.topic"
	err = ch.ExchangeDeclare(
		exName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed declare Exchange")
	}

	q, err := ch.QueueDeclare(
		"bark-bridge",
		true,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed declare queue")
	}

	err = ch.QueueBind(
		q.Name,             // queue name
		"*.iot.sms.upload", // routing key
		exName,             // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed bind queue")
	}

	msgCh, err := ch.Consume(
		q.Name,
		"bark-bridge", // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed create consume chan")
	}

	p.wg.Add(1)
	go func() {

		for {
			select {
			case msg := <-msgCh:
				logrus.Info(string(msg.Body))
				// p.ProvideCh <- barkserver.Message{}
			case _, ok := <-p.stopCh:
				if !ok {
					ch.Close()
					conn.Close()
					close(p.ProvideCh)
					p.wg.Done()
					logrus.Info(fmt.Sprintf("exit %s", p.GetName()))
					return
				}
			}
		}

	}()

	logrus.Info(fmt.Sprintf("start %s", p.GetName()))
}

func (p *AmqpProvider) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *AmqpProvider) GetCh() <-chan barkserver.Message {
	return p.ProvideCh
}

func (p *AmqpProvider) GetName() string {
	return "AmqpProvider"
}
