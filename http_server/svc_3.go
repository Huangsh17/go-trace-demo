package http_server

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go-trace-demo/grpc_client"
	"go-trace-demo/middleware"
	"net/http"
)

func Router03() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(middleware.Trace())
	e.GET("/r3", func(c *gin.Context) {
		value, _ := c.Get("span")
		span, _ := value.(opentracing.Span)
		defer span.Finish()
		valueTracer, _ := c.Get("tracer")
		tracer, _ := valueTracer.(opentracing.Tracer)
		grpc_client.GrpcClient(tracer, span)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Welcome server 03",
			},
		)
	})
	return e
}
