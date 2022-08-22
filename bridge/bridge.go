package bridge

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/c1emon/barkbridge/barkserver"
	"github.com/sirupsen/logrus"
)

type Provider interface {
	Start()
	Stop()
	GetCh() <-chan barkserver.Message
	GetName() string
}

type Bridge struct {
	Providers   map[string]Provider
	osSigs      chan os.Signal
	msgCh       chan barkserver.Message
	wg          *sync.WaitGroup
	BarkAddress string
}

func New(server string) *Bridge {
	b := &Bridge{
		Providers:   make(map[string]Provider),
		osSigs:      make(chan os.Signal, 1),
		msgCh:       make(chan barkserver.Message),
		wg:          &sync.WaitGroup{},
		BarkAddress: server,
	}

	return b
}

func (b *Bridge) AddProvider(id string, p Provider) {
	logrus.WithFields(logrus.Fields{
		"id":   id,
		"type": p.GetName(),
	}).Info("add provider")
	b.Providers[id] = p
}

func (b *Bridge) Server() {
	signal.Notify(b.osSigs, syscall.SIGINT, syscall.SIGTERM)

	for id, provider := range b.Providers {
		b.wg.Add(1)
		go func(id string, p Provider) {
			// logrus.WithField("id", id).Info("start bridge provider")
			for msg := range p.GetCh() {
				b.msgCh <- msg
			}
			b.wg.Done()
			// logrus.WithField("id", id).Info("stop bridge provider")
		}(id, provider)

	}

	go func() {
		for msg := range b.msgCh {
			logrus.WithFields(logrus.Fields{
				"id":    msg.Id,
				"title": msg.Title,
			}).Info("bridge message")
			barkserver.Push(b.BarkAddress, msg)
		}
	}()

	for id, provider := range b.Providers {
		logrus.WithFields(logrus.Fields{
			"id":   id,
			"type": provider.GetName(),
		}).Info("start provider")
		provider.Start()
	}

	logrus.Info("bridge start!")

	<-b.osSigs
	logrus.Debug("wait for stop")
	for _, provider := range b.Providers {
		go func(provider Provider) {
			provider.Stop()
		}(provider)
	}
	b.wg.Wait()
	close(b.msgCh)
	logrus.Info("bridge exit!")

}
