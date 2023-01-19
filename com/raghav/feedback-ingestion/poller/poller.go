package poller

import (
	"com/raghav/feedback-ingestion/feedback"
	"time"
)

type Poller struct {
	serviceFactory feedback.ServiceFactory
	feedbackStore  feedback.Store
	ticker         *time.Ticker
	stop           chan bool
}

func NewPoller(serviceFactory feedback.ServiceFactory, feedbackStore feedback.Store) *Poller {
	return &Poller{
		serviceFactory: serviceFactory,
		feedbackStore:  feedbackStore,
		ticker:         time.NewTicker(15 * time.Minute),
		stop:           make(chan bool),
	}
}

func (p *Poller) Start() {
	for {
		select {
		case <-p.ticker.C:
			services := p.serviceFactory.GetAll()
			for _, service := range services {
				data, err := service.RetrieveData()
				if err != nil {
					continue
				}
				transformedData, err := service.TransformData(data)
				if err != nil {
					continue
				}
				p.feedbackStore.StoreData(transformedData)
			}
		case <-p.stop:
			p.ticker.Stop()
			return
		}
	}
}

func (p *Poller) Stop() {
	p.stop <- true
}
