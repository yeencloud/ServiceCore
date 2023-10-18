package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type RequestMetadata struct {
	StartTime time.Time
	IP        string
	Method    string
	Path      string
}

func (r RequestMetadata) MarshalZerologObject(e *zerolog.Event) {
	e.Time("at", r.StartTime)
	e.Str("ip", r.IP)
	e.Str("path", r.Path)
	e.Str("method", r.Method)
}

func (shs *ServiceHTTPServer) logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		ri := RequestMetadata{
			IP:        clientIP,
			StartTime: startTime,
			Method:    method,
			Path:      path,
		}

		c.Set("requestmetadata", ri)
		c.Next()
		endTime := time.Now()

		code := c.Writer.Status()

		logWeight := log.Info()
		if code >= 300 {
			logWeight = log.Warn()
		}

		currentLog := logWeight.Object("request", ri).
			Int("code", code).
			TimeDiff("duration", endTime, startTime)

		currentLog.Msg("request served")
	}
}
