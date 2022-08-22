package providers

import (
	"sync"

	"github.com/c1emon/barkbridge/barkserver"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpProvider struct {
	Addr      string
	ProvideCh chan barkserver.Message
	stopCh    chan any
	wg        *sync.WaitGroup
}

func (p *AmqpProvider) Start() {
	amqp.Dial(p.Addr)
}

func (p *AmqpProvider) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *AmqpProvider) GetCh() <-chan barkserver.Message {

	return p.ProvideCh
}

func (p *AmqpProvider) GetName() string {
	return ""
}
