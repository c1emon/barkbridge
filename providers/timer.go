package providers

import (
	"fmt"
	"time"

	"github.com/c1emon/barkbridge/barkserver"
	"github.com/c1emon/barkbridge/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type TimeProvider struct {
	ProvideCh chan barkserver.Message
	name      string
	cronExp   string
	body      string
	title     string
	devicekey string
	cron      *cron.Cron
}

func NewTimeProvider(body, title, devicekey, cron string) *TimeProvider {
	p := &TimeProvider{
		ProvideCh: make(chan barkserver.Message),
		name:      fmt.Sprintf("timer-provider-%s", utils.RandStr(5)),
		cronExp:   cron,
		body:      body,
		title:     title,
		devicekey: devicekey,
	}

	return p
}

func (p *TimeProvider) prepare() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc(p.cronExp, func() {

		logrus.WithField("name", p.name).Debug("timer provider")
		p.ProvideCh <- barkserver.Message{
			Title:     p.title,
			Body:      p.title,
			DeviceKey: p.devicekey,
			Category:  "category",
			Id:        time.Now().Format("2006-01-02 15:04:05"),
		}
	})

	p.cron = c
}

func (p *TimeProvider) Start() {
	p.prepare()

	p.cron.Start()
	logrus.Info(fmt.Sprintf("start %s", p.GetName()))
}

func (p *TimeProvider) Stop() {
	p.cron.Stop()
	close(p.ProvideCh)

	logrus.Info(fmt.Sprintf("exit %s", p.GetName()))
}

func (p *TimeProvider) GetCh() <-chan barkserver.Message {
	return p.ProvideCh
}

func (p *TimeProvider) GetName() string {
	return p.name
}
