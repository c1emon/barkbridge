package providers

import (
	"fmt"
	"sync"

	"github.com/c1emon/barkbridge/barkserver"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type TimeProvider struct {
	ProvideCh chan barkserver.Message
	stopCh    chan any
	wg        *sync.WaitGroup
}

func NewTimeProvider() *TimeProvider {
	p := &TimeProvider{
		wg:        &sync.WaitGroup{},
		ProvideCh: make(chan barkserver.Message),
		stopCh:    make(chan any),
	}

	return p
}

func (p *TimeProvider) Start() {

	p.wg.Add(1)
	go func() {
		i := 0
		c := cron.New(cron.WithSeconds())

		c.AddFunc("*/2 * * * * ?", func() {
			i++
			logrus.WithField("idx", i).Debug("timer provider")
			p.ProvideCh <- barkserver.Message{
				Title:     fmt.Sprintf("Test: %d", i),
				Body:      "test from timer provider",
				DeviceKey: "key",
				Category:  "category",
				Id:        fmt.Sprintf("%d", i),
			}
		})

		c.Start()

		for _ = range p.stopCh {
		}
		c.Stop()
		close(p.ProvideCh)
		p.wg.Done()
		logrus.Info(fmt.Sprintf("exit %s", p.GetName()))
	}()
	logrus.Info(fmt.Sprintf("start %s", p.GetName()))
}

func (p *TimeProvider) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *TimeProvider) GetCh() <-chan barkserver.Message {

	return p.ProvideCh
}

func (p *TimeProvider) GetName() string {
	return "TimeProvider"
}
