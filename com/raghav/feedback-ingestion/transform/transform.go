package transform

package transform

import (
    "com/raghav/feedback-ingestion/model"
)

// Transformer interface defines the TransformData method that should be implemented by each source
type Transformer interface {
    TransformData(data []model.FeedbackData) ([]model.Feedback, error)
}

// DiscourseTransformer struct implements the Transformer interface for the Discourse source
type DiscourseTransformer struct{}

// TransformData method transforms the input data specific to the Discourse source
func (d *DiscourseTransformer) TransformData(data []model.FeedbackData) ([]model.Feedback, error) {
    var feedbacks []model.Feedback
    for _, v := range data {
        feedbacks = append(feedbacks, model.Feedback{
            ID:         v.ID,
            TopicID:    v.TopicID,
            CreatedAt:  v.CreatedAt,
            LikeCount:  v.LikeCount,
            Blurb:      v.Blurb,
            PostNumber: v.PostNumber,
            TopicTitle: v.TopicTitle,
        })
    }
    return feedbacks, nil
}


