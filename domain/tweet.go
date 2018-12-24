package domain

import (
	"context"
	"encoding/json"
	"time"
)

// Tweet represent the tweet data' struct
type Tweet struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	CreatedTime time.Time `json:"createdTime"`
}

func (t *Tweet) ToJSON() (string, error) {
	byt, err := json.Marshal(t)
	return string(byt), err
}

func (t *Tweet) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), &t)
}

// TweetQueryParam represent the fetch tweet query param
type TweetQueryParam struct {
	Cursor string
	Num    int
}

// TweetUsecase represent the tweet usecase
type TweetUsecase interface {
	Post(ctx context.Context, t *Tweet) error
	Fetch(ctx context.Context, qparam TweetQueryParam) (res []Tweet, cursor string, err error)
	Get(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}

// TweetRepository represent the tweet repository
type TweetRepository interface {
	Create(ctx context.Context, t *Tweet) error
	Fetch(ctx context.Context, qparam TweetQueryParam) (res []Tweet, cursor string, err error)
	Get(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}
