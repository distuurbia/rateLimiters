package main

import (
	"log"
	"net/http"
	"time"

	customMiddlewarre "github.com/distuurbia/rateLimiters/examples/middleware"
	"github.com/distuurbia/rateLimiters/fixedWindow"
	"github.com/distuurbia/rateLimiters/slidingLog"
	"github.com/distuurbia/rateLimiters/slidingWindow"
	"github.com/distuurbia/rateLimiters/tokenBucket"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return client
}

type helloWorker struct{}

func newServerWorker() *helloWorker {
	return &helloWorker{}
}

func (sw *helloWorker) Hello(c echo.Context) error {
	return c.JSON(http.StatusAccepted, "HELLO")
}

func main() {
	const (
		intervalToBeParsed = "10s"
		maxRequests        = 2
		maxTokens          = 500
		rate               = 50
	)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := newServerWorker()
	interval, err := time.ParseDuration(intervalToBeParsed)
	if err != nil {
		log.Fatal(err)
	}

	client := connectRedis()
	fw := fixedWindow.NewFixedWindow(interval, maxRequests)
	tb := tokenBucket.NewTokenBucket(rate, maxTokens)
	sl := slidingLog.NewSlidingLog(client)
	sw := slidingWindow.NewSlidingWindow(client)
	cm := customMiddlewarre.NewCustomMiddleware(fw, tb, sl, sw)
	e.GET("/helloFW", h.Hello, cm.MiddlewareFixedWindow)
	e.GET("/helloTB", h.Hello, cm.MiddlewareTokenBucket)
	e.GET("/helloSL", h.Hello, cm.MiddlewareSlidingLog)
	e.GET("/helloSW", h.Hello, cm.MiddlewareSlidingWindow)
	e.Logger.Fatal(e.Start(":8080"))
}
