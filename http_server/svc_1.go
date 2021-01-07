package http_server

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go-trace-demo/http_client"
	"go-trace-demo/middleware"
	"net/http"
)

func Router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(middleware.Trace())
	e.GET("/r1", func(c *gin.Context) {
		valueTracer, _ := c.Get("tracer")
		tracer, _ := valueTracer.(opentracing.Tracer)
		value, _ := c.Get("span")
		span, _ := value.(opentracing.Span)
		http_client.HttpClient("http://127.0.0.1:8082/r2", tracer, span)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Welcome server 01",
			},
		)
	})
	return e
}
