package providers

import (
	"fmt"
	"sync"
	"time"

	"github.com/c1emon/barkbridge/barkserver"
	"github.com/c1emon/barkbridge/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type TimeProvider struct {
	ProvideCh chan barkserver.Message
	stopCh    chan any
	wg        *sync.WaitGroup
	name      string
	cron      string
	body      string
	title     string
	devicekey string
}

func NewTimeProvider(body, title, devicekey, cron string) *TimeProvider {
	p := &TimeProvider{
		wg:        &sync.WaitGroup{},
		ProvideCh: make(chan barkserver.Message),
		stopCh:    make(chan any),
		name:      fmt.Sprintf("timer-provider-%s", utils.RandStr(5)),
		cron:      cron,
		body:      body,
		title:     title,
		devicekey: devicekey,
	}

	return p
}

func (p *TimeProvider) Start() {

	p.wg.Add(1)
	go func() {

		c := cron.New(cron.WithSeconds())

		c.AddFunc(p.cron, func() {

			logrus.WithField("name", p.name).Debug("timer provider")
			p.ProvideCh <- barkserver.Message{
				Title:     p.title,
				Body:      p.title,
				DeviceKey: p.devicekey,
				Category:  "category",
				Id:        time.Now().Format("2006-01-02 15:04:05"),
			}
		})

		c.Start()

		for range p.stopCh {
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
	return p.name
}
