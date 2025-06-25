package logger

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var body string
		if logrus.IsLevelEnabled(logrus.DebugLevel) && c.Request.Body != nil {
			buf, err := io.ReadAll(c.Request.Body)
			if err == nil {
				body = string(buf)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
			}
		}

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logFields := logrus.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"referer":    c.Request.Referer(),
			"duration":   duration.String(),
			"user_agent": c.Request.UserAgent(),
		})

		if logrus.IsLevelEnabled(logrus.DebugLevel) && body != "" {
			logFields = logFields.WithField("body", body)
		}

		if status >= 500 {
			logFields.Error("internal server error")
		} else {
			logFields.Info("handled request")
		}
	}
}
