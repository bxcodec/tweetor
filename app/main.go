package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	handler "github.com/bxcodec/tweetor/tweet/delivery/http"
	"github.com/bxcodec/tweetor/tweet/repository"
	"github.com/bxcodec/tweetor/tweet/usecase"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.toml")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		// Log Error
		log.Fatal(err)
	}
}

func main() {
	redisAddr := viper.GetString("redis.address")
	redisDB := viper.GetInt("redis.db")
	redisPass := viper.GetString("redis.pass")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass, // no password set
		DB:       redisDB,   // use default DB
	})
	ctxDuration := viper.GetDuration("context.timeout")
	repo := repository.NewRedisRepository(redisClient)
	ucase := usecase.New(repo, time.Second*ctxDuration)

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong!")
	})
	handler.AddTweetHandler(e, ucase)

	logrus.Error(e.Start(viper.GetString("address")))
}
