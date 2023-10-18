package middleware

import (
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLoggerMiddleware() echo.MiddlewareFunc {
	log.Out = os.Stdout
	log.SetLevel(logrus.InfoLevel)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			log.WithFields(logrus.Fields{
				"method":   c.Request().Method,
				"url":      c.Request().Response.Request.URL,
				"status":   c.Response().Status,
				"latency":  latency,
				"clientIP": c.RealIP(),
			}).Info("Request handled")

			return err
		}
	}
}
