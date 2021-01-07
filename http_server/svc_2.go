package http_server

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go-trace-demo/http_client"
	"go-trace-demo/middleware"
	"net/http"
)

func Router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(middleware.Trace())
	e.GET("/r2", func(c *gin.Context) {
		value, _ := c.Get("span")
		span, _ := value.(opentracing.Span)
		valueTracer, _ := c.Get("tracer")
		tracer, _ := valueTracer.(opentracing.Tracer)
		http_client.HttpClient("http://127.0.0.1:8083/r3", tracer, span)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Welcome server 02",
			},
		)
	})
	return e
}
