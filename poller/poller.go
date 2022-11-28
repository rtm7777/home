package poller

import (
	"time"
)

type Poller struct {
	Tick         chan time.Time
	ticker       *time.Ticker
	timer        time.Duration
	isThrottling bool
}

func NewPoller(timer time.Duration) Poller {
	return Poller{
		Tick:         make(chan time.Time),
		ticker:       time.NewTicker(timer),
		timer:        timer,
		isThrottling: false,
	}
}

func (p *Poller) Init() {
	for {
		select {
		case t := <-p.ticker.C:
			p.Tick <- t
		}
	}
}

func (p *Poller) Throttle(throttleState bool) {
	if p.isThrottling != throttleState {
		if throttleState {
			p.ticker = time.NewTicker(p.timer * 10)
		} else {
			p.ticker = time.NewTicker(p.timer)
		}
	}
}
