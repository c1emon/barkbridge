package bridge

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/c1emon/barkbridge/barkserver"
)

type Provider interface {
	Start()
	Stop()
	GetCh() <-chan barkserver.Message
}

type Bridge struct {
	Providers map[string]Provider
	osSigs    chan os.Signal
	msgCh     chan barkserver.Message
	wg        *sync.WaitGroup
}

func New() *Bridge {
	b := &Bridge{
		Providers: nil,
		osSigs:    make(chan os.Signal, 1),
		msgCh:     make(chan barkserver.Message),
		wg:        &sync.WaitGroup{},
	}

	return b
}

func (b *Bridge) Server() {
	signal.Notify(b.osSigs, syscall.SIGINT, syscall.SIGTERM)

	for _, provider := range b.Providers {
		b.wg.Add(1)
		go func(p Provider) {
			for msg := range p.GetCh() {
				b.msgCh <- msg
			}
			b.wg.Done()
		}(provider)

	}

	b.wg.Add(1)
	go func() {
		for msg := range b.msgCh {
			barkserver.Push("", msg)
		}
		b.wg.Done()
	}()

	<-b.osSigs
	for _, provider := range b.Providers {
		provider.Stop()
	}
	close(b.msgCh)
	b.wg.Wait()

}
