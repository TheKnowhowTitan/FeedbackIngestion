package feedback

import (
	"api"
	"com/raghav/feedback-ingestion/service"
	"com/raghav/feedback-ingestion/store"
	"com/raghav/feedback-ingestion/transform"
	"fmt"
	"time"
)

type Feedback struct {
	ServiceRegistry *service.ServiceRegistry
	Store           *store.FeedbackStore
	Poller          *api.Poller
}

// NewFeedback returns a new instance of feedback
func NewFeedback() *Feedback {
	serviceRegistry := service.NewServiceRegistry()
	store := store.NewFeedbackStore()
	poller := api.NewPoller(serviceRegistry, store)
	go poller.Start()

	return &Feedback{
		ServiceRegistry: serviceRegistry,
		Store:           store,
		Poller:          poller,
	}
}

// AddFeedbackSource adds a new feedback source to the system
func (f *Feedback) AddFeedbackSource(name string, client *transform.Client, transformFunc transform.TransformFunc, pollInterval time.Duration) error {
	if err := f.ServiceRegistry.RegisterService(name, client, transformFunc); err != nil {
		return fmt.Errorf("error while registering service: %v", err)
	}

	f.Poller.AddSource(name, pollInterval)
	return nil
}

// GetFeedback returns the feedback for the given source and time range
func (f *Feedback) GetFeedback(source string, from, to time.Time) ([]store.Feedback, error) {
	return f.Store.Get(source, from, to)
}
