package providers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/c1emon/barkbridge/barkserver"
	"github.com/c1emon/barkbridge/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AmqpProvider struct {
	ProvideCh chan barkserver.Message
	stopCh    chan any
	wg        *sync.WaitGroup
	name      string
	conf      *AmqpConf
	conn      *amqp.Connection
	ch        *amqp.Channel
	msgCh     <-chan amqp.Delivery
}

type AmqpConf struct {
	Addr       string
	Exchange   string
	Queue      string
	Topic      string
	RoutingKey string
}

func NewAmqpProvider(conf *AmqpConf) *AmqpProvider {
	p := &AmqpProvider{
		wg:        &sync.WaitGroup{},
		ProvideCh: make(chan barkserver.Message),
		stopCh:    make(chan any),
		conf:      conf,
		name:      fmt.Sprintf("amqp-provider-%s", utils.RandStr(5)),
	}

	return p
}

func (p *AmqpProvider) prepare() {
	conn, err := amqp.Dial(p.conf.Addr)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed dial amqp server")
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.WithField("error", err).Fatal("failed open channel")
	}

	err = ch.ExchangeDeclare(
		p.conf.Exchange,
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
		p.conf.Queue,
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
		q.Name,            // queue name
		p.conf.RoutingKey, // routing key
		p.conf.Exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed bind queue")
	}

	msgCh, err := ch.Consume(
		q.Name,
		fmt.Sprintf("bark-bridge-%s", utils.RandStr(5)), // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed create consume chan")
	}

	p.conn = conn
	p.ch = ch
	p.msgCh = msgCh

}

func (p *AmqpProvider) Start() {
	p.prepare()

	p.wg.Add(1)
	go func() {
		defer p.conn.Close()
		defer p.ch.Close()
		for {
			select {
			case msg := <-p.msgCh:
				logrus.Debug(fmt.Sprintf("amqp msg:\n%s", string(msg.Body)))
				m := barkserver.Message{}
				err := json.Unmarshal(msg.Body, &m)
				if err != nil {
					logrus.Warn(fmt.Sprintf("unmarshal err: %s", err))
					continue
				}
				p.ProvideCh <- m
			case _, ok := <-p.stopCh:
				if !ok {
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
	return p.name
}
