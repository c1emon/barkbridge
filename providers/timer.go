package providers

import (
	"fmt"
	"sync"
	"time"

	"github.com/c1emon/barkbridge/barkserver"
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
		for {

			select {
			case _, ok := <-p.stopCh:
				if !ok {
					close(p.ProvideCh)
					p.wg.Done()
					return
				}
			default:
				time.Sleep(time.Second * time.Duration(5))
				i++
				logrus.WithField("idx", i).Debug("timer provider")
				p.ProvideCh <- barkserver.Message{
					Title:     fmt.Sprintf("Test: %d", i),
					Body:      "test from timer provider",
					DeviceKey: "key",
					Category:  "category",
					Id:        fmt.Sprintf("%d", i),
				}
			}

		}
	}()
	logrus.Info("timer provider started")
}

func (p *TimeProvider) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *TimeProvider) GetCh() <-chan barkserver.Message {

	return p.ProvideCh
}
