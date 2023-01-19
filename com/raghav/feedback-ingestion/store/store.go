package store

import (
	"com/raghav/feedback-ingestion/feedback"
)

type Store struct {
	feedbackData map[string][]feedback.Feedback
}

func NewStore() *Store {
	return &Store{
		feedbackData: make(map[string][]feedback.Feedback),
	}
}

func (s *Store) AddFeedback(source string, feedbacks []feedback.Feedback) {
	s.feedbackData[source] = feedbacks
}

func (s *Store) GetFeedback(source string) ([]feedback.Feedback, bool) {
	feedbacks, ok := s.feedbackData[source]
	return feedbacks, ok
}
