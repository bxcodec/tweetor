package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bxcodec/tweetor/domain"
	"github.com/go-redis/redis"
)

type redisHandler struct {
	client *redis.Client
}

var (
	duration  = time.Second * 300
	prefixKey = "tweet"
)

func NewRedisRepository(client *redis.Client) domain.TweetRepository {
	return &redisHandler{
		client: client,
	}
}

func (r redisHandler) Create(ctx context.Context, t *domain.Tweet) error {
	key := fmt.Sprintf("%s:%s", prefixKey, t.ID)
	jsonStr, err := t.ToJSON()
	if err != nil {
		return err
	}
	err = r.client.Set(key, jsonStr, duration).Err()
	if err != nil {
		return err
	}
	err = r.client.SAdd(prefixKey, key).Err()
	return err
}

func (r redisHandler) Fetch(ctx context.Context, qparam domain.TweetQueryParam) (res []domain.Tweet, cursor string, err error) {
	arrKeys := []string{}
	var decodedCursor int64
	if qparam.Cursor != "" {
		decodedCursor, err = decodeCursor(qparam.Cursor)
		if err != nil {
			return res, cursor, err
		}
	}

	status := r.client.Sort(
		prefixKey,
		&redis.Sort{
			Alpha:  true,
			Count:  int64(qparam.Num),
			Order:  "desc",
			Offset: decodedCursor,
		},
	)

	err = status.ScanSlice(&arrKeys)
	if err != nil {
		return
	}

	res = []domain.Tweet{}
	for _, key := range arrKeys {
		item, err := r.Get(ctx, key)
		if err != nil {
			return res, cursor, err
		}
		res = append(res, item)

	}
	if len(res) > 0 {
		cursor = encodeCursor(fmt.Sprintf("%d", decodedCursor+int64(qparam.Num)))
	}
	return
}

func (r redisHandler) Get(ctx context.Context, id string) (domain.Tweet, error) {
	res := domain.Tweet{}
	status := r.client.Get(fmt.Sprintf("%s:%s", prefixKey, id))
	jbyt, err := status.Bytes()
	if err != nil {
		if err == redis.Nil {
			return res, domain.ErrNotFound
		}
		return res, err
	}

	err = json.Unmarshal(jbyt, &res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (r redisHandler) Delete(ctx context.Context, id string) error {
	inctCmd := r.client.Unlink(fmt.Sprintf("%s:%s", prefixKey, id))
	err := inctCmd.Err()

	if err == redis.Nil {
		return domain.ErrNotFound
	}
	return err

}

func decodeCursor(encodedTime string) (int64, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return 0, err
	}
	idString := string(byt)

	res, err := strconv.ParseInt(idString, 10, 64)

	return res, err
}

func encodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}
