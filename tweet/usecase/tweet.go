package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/tweetor/domain"
)

type usecase struct {
	repository domain.TweetRepository
	timeout    time.Duration
}

func New(repo domain.TweetRepository, timeout time.Duration) domain.TweetUsecase {
	return &usecase{
		repository: repo,
		timeout:    timeout,
	}
}

func (u usecase) Post(c context.Context, t *domain.Tweet) (err error) {
	if c == nil {
		err = domain.ErrContextNil
		return
	}
	t.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	t.CreatedTime = time.Now()
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	err = u.repository.Create(ctx, t)
	return
}

func (u usecase) Fetch(c context.Context, qparam domain.TweetQueryParam) (res []domain.Tweet, cursor string, err error) {
	if c == nil {
		err = domain.ErrContextNil
		return
	}
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	res, cursor, err = u.repository.Fetch(ctx, qparam)

	return
}

func (u usecase) Get(c context.Context, id string) (res domain.Tweet, err error) {
	if c == nil {
		err = domain.ErrContextNil
		return
	}
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	res, err = u.repository.Get(ctx, id)
	return
}
func (u usecase) Delete(c context.Context, id string) (err error) {
	if c == nil {
		err = domain.ErrContextNil
		return
	}
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()
	err = u.repository.Delete(ctx, id)
	return
}
