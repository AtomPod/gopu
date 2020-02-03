package middleware

import (
	"context"
	"time"

	"github.com/ngs24313/gopu/utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

type levelFunc func(msg string, fields ...zapcore.Field)

var (
	levelMapping = map[zapcore.Level]func(context.Context) levelFunc{
		zapcore.DebugLevel: func(c context.Context) levelFunc {
			return log.Logger(c).Debug
		},
		zapcore.ErrorLevel: func(c context.Context) levelFunc {
			return log.Logger(c).Error
		},
		zapcore.InfoLevel: func(c context.Context) levelFunc {
			return log.Logger(c).Info
		},
		zapcore.WarnLevel: func(c context.Context) levelFunc {
			return log.Logger(c).Warn
		},
		zapcore.FatalLevel: func(c context.Context) levelFunc {
			return log.Logger(c).Fatal
		},
		zap.DPanicLevel: func(c context.Context) levelFunc {
			return log.Logger(c).DPanic
		},
	}
)

//Logger gin log middleware
func Logger(lv ...zapcore.Level) gin.HandlerFunc {

	logCaller := levelMapping[zap.DebugLevel]
	if len(lv) > 0 {
		if caller, ok := levelMapping[lv[0]]; ok {
			logCaller = caller
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		params := gin.LogFormatterParams{
			Request:      c.Request,
			Keys:         c.Keys,
			StatusCode:   c.Writer.Status(),
			Latency:      time.Since(start),
			Method:       c.Request.Method,
			ClientIP:     c.ClientIP(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     c.Writer.Size(),
		}

		if raw != "" {
			path = path + "?" + raw
		}
		params.Path = path

		logCaller(c.Request.Context())(
			"Receive a http request",
			zap.String("path", params.Path),
			zap.String("method", params.Method),
			zap.Int("status_code", params.StatusCode),
			zap.Duration("latency", params.Latency),
			zap.String("client_ip", params.ClientIP),
			zap.String("error", params.ErrorMessage),
			zap.Int("body_size", params.BodySize),
			zap.Any("keys", params.Keys),
		)
	}
}
