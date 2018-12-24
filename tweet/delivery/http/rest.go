package http

import (
	"net/http"
	"strconv"

	"github.com/bxcodec/tweetor/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type handlerTweet struct {
	usecase domain.TweetUsecase
}

func AddTweetHandler(e *echo.Echo, usecase domain.TweetUsecase) {
	handler := handlerTweet{
		usecase: usecase,
	}
	e.POST("/tweets", handler.Post)
	e.GET("/tweets", handler.Fetch)
	e.GET("/tweets/:id", handler.Get)
	e.DELETE("/tweets/:id", handler.Delete)
}

func (h *handlerTweet) Post(c echo.Context) error {

	item := domain.Tweet{}

	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	err = h.usecase.Post(c.Request().Context(), &item)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, item)
}

func (h *handlerTweet) Get(c echo.Context) error {
	id := c.Param("id")
	res, err := h.usecase.Get(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *handlerTweet) Fetch(c echo.Context) (err error) {
	cursor := c.QueryParam("cursor")
	numStr := c.QueryParam("limit")
	num := 10

	if numStr != "" {
		num, err = strconv.Atoi(numStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	}
	qparam := domain.TweetQueryParam{
		Num:    num,
		Cursor: cursor,
	}
	res, nextCursor, err := h.usecase.Fetch(c.Request().Context(), qparam)
	if err != nil {
		if err == domain.ErrNotFound {
			return c.JSON(http.StatusOK, res)
		}
		logrus.Error(err)
		return err
	}

	c.Response().Header().Set("X-Cursor", nextCursor)
	return c.JSON(http.StatusOK, res)
}

func (h *handlerTweet) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.usecase.Delete(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		logrus.Error(err)
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
