package discourse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Discourse struct {
	BaseUrl string
	Client  *http.Client
}

type DiscourseData struct {
	Posts []struct {
		ID         int    `json:"id"`
		TopicID    int    `json:"topic_id"`
		CreatedAt  string `json:"created_at"`
		LikeCount  int    `json:"like_count"`
		PostNumber int    `json:"post_number"`
		TopicTitle string `json:"topic_title"`
		Blurb      string `json:"blurb"`
	} `json:"posts"`
}

type DiscoursePostData struct {
	PostStream struct {
		Posts []struct {
			ID         int    `json:"id"`
			TopicID    int    `json:"topic_id"`
			CreatedAt  string `json:"created_at"`
			Cooked     string `json:"cooked"`
			PostNumber int    `json:"post_number"`
			TopicTitle string `json:"topic_title"`
		} `json:"posts"`
	} `json:"post_stream"`
}

func NewDiscourse(baseUrl string) *Discourse {
	return &Discourse{
		BaseUrl: baseUrl,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (d *Discourse) RetrieveData(startTime, endTime time.Time) ([]DiscourseData, error) {
	var data []DiscourseData
	// Get list of posts in the given time range
	params := fmt.Sprintf("search.json?page=1&q=after:%s+before:%s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	resp, err := d.Client.Get(d.BaseUrl + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal the JSON response
	var searchResp struct {
		Posts []struct {
			ID        int    `json:"id"`
			TopicID   int    `json:"topic_id"`
			PostURL   string `json:"url"`
			CreatedAt string `json:"created_at"`
		} `json:"posts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	// Get full post data for each post
	for _, post := range searchResp.Posts {
		postData, err := d.fetchFullPostData(post.ID, post.TopicID)
		if err != nil {
			return nil, err
		}
		data = append(data, postData)
	}

	return data, nil
}

func (d *Discourse) fetchFullPostData(postID, topicID int) (DiscourseData, error) {
	url := fmt.Sprintf("%st/%d/posts.json?post_ids[]=%d", d.BaseUrl, topicID, postID)
	resp, err := http.Get(url)
	if err != nil {
		return DiscourseData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return DiscourseData{}, fmt.Errorf("failed to fetch post data, status code: %d", resp.StatusCode)
	}

	var postResponse struct {
		PostStream struct {
			Posts []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Username  string `json:"username"`
				CreatedAt string `json:"created_at"`
				Cooked    string `json:"cooked"`
				TopicID   int    `json:"topic_id"`
				TopicSlug string `json:"topic_slug"`
			} `json:"posts"`
		} `json:"post_stream"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&postResponse); err != nil {
		return DiscourseData{}, err
	}

	if len(postResponse.PostStream.Posts) == 0 {
		return DiscourseData{}, fmt.Errorf("no post found for postID: %d and topicID: %d", postID, topicID)
	}

	post := postResponse.PostStream.Posts[0]
	createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
	if err != nil {
		return DiscourseData{}, err
	}

	return DiscourseData{
		ID:        post.ID,
		Name:      post.Name,
		Username:  post.Username,
		CreatedAt: createdAt,
		Cooked:    post.Cooked,
		TopicID:   post.TopicID,
		TopicSlug: post.TopicSlug,
	}, nil
}
