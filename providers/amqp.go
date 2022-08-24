package providers

import (
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
		logrus.WithField("addr", p.Addr).Panic("failed dial amqp server")
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.WithField("addr", p.Addr).Fatal("failed open channel")
	}

	exName := "mqtt"
	err = ch.ExchangeDeclare(
		exName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("exchange", exName).Fatal("failed connect exchaneg")
	}

	q, err := ch.QueueDeclare(
		"name",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("queue", "").Fatal("")
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		exName, // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("", "").Fatal("")
	}

	msgCh, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logrus.WithField("", "").Fatal("")
	}

	p.wg.Add(1)
	go func() {

		for {
			select {
			case msg := <-msgCh:
				logrus.Info(msg)
				p.ProvideCh <- barkserver.Message{}
			case _, ok := <-p.stopCh:
				if !ok {
					ch.Close()
					conn.Close()
					close(p.ProvideCh)
					p.wg.Done()
					logrus.Info("exit")
					return
				}
			}
		}

	}()

	logrus.Info("start provider")
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
