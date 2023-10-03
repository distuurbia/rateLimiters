package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/distuurbia/rateLimiters/fixedWindow"
	"github.com/distuurbia/rateLimiters/slidingLog"
	"github.com/distuurbia/rateLimiters/slidingWindow"
	"github.com/distuurbia/rateLimiters/tokenBucket"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type CustomMiddleware struct {
	fw *fixedWindow.FixedWindow
	tb *tokenBucket.TokenBucket
	sl *slidingLog.SlidingLog
	sw *slidingWindow.SlidingWindow
}

func NewCustomMiddleware(
	fw *fixedWindow.FixedWindow,
	tb *tokenBucket.TokenBucket,
	sl *slidingLog.SlidingLog,
	sw *slidingWindow.SlidingWindow) *CustomMiddleware {
	return &CustomMiddleware{
		fw: fw,
		tb: tb,
		sl: sl,
		sw: sw,
	}
}

func (cm *CustomMiddleware) MiddlewareFixedWindow(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := cm.fw.CheckIfRequestAllowed()
		if !res {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later fixedWindow")
		}
		return next(c)
	}
}

func (cm *CustomMiddleware) MiddlewareTokenBucket(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const tokens = 100
		res := cm.tb.CheckIfRequestAllowed(tokens)
		if !res {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later tokenBucket")
		}
		return next(c)
	}
}

func (cm *CustomMiddleware) MiddlewareSlidingLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const (
			maxRequests        = 10
			intervalToBeParsed = "10s"
		)
		interval, err := time.ParseDuration(intervalToBeParsed)
		if err != nil {
			err = fmt.Errorf("CustomMiddleware-time.ParseDuration-err: %w", err)
			slog.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		res := cm.sl.CheckIfRequestAllowed(c.Request().Context(), c.RealIP(), uuid.NewString(), interval, maxRequests)
		if !res {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later sliding log")
		}
		return next(c)
	}
}

func (cm *CustomMiddleware) MiddlewareSlidingWindow(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const (
			maxRequests        = 10
			intervalToBeParsed = "10s"
		)
		interval, err := time.ParseDuration(intervalToBeParsed)
		if err != nil {
			err = fmt.Errorf("CustomMiddleware-time.ParseDuration-err: %w", err)
			slog.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		res, err := cm.sw.CheckIfRequestAllowed(c.Request().Context(), c.RealIP(), interval, maxRequests)
		if err != nil {
			err = fmt.Errorf("CustomMiddleware-cm.sw.CheckIfRequestAllowed-err: %w", err)
			slog.Error(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if !res {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later sliding window")
		}
		return next(c)
	}
}
