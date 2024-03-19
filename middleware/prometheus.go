package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// api request
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests received.",
		},
		[]string{"method", "URL"},
	)

	// api response time
	httpRequestsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "URL"},
	)

	httpRequestsSuccess = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_success_total",
			Help: "Total number of requests errors.",
		},
		[]string{"method", "URL", "status"}, // 使用路径作为标签
	)
)

func PrometheusRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		c.Next()
	}
}

func PrometheusDuration() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		httpRequestsDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(duration.Milliseconds()))
	}
}

func PrometheusResErr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		statusCode := c.Writer.Status()
		httpRequestsSuccess.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(statusCode)).Inc()
	}
}
