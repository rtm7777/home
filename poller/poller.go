package poller

import (
	"time"
)

type poller struct {
	ticker *time.Ticker
}

func NewPoller(timer time.Duration) *poller {
	return &poller{
		ticker: time.NewTicker(timer),
	}
}

func (p *poller) Run(tick chan<- time.Time) {
	for {
		select {
		case t := <-p.ticker.C:
			tick <- t
		}
	}
}
