package middleware

import (
	"github.com/gin-gonic/gin"
	"go-trace-demo/tracing"
)

func Trace() gin.HandlerFunc {
	return func(context *gin.Context) {
		tracer, closer, _ := tracing.InitTracer("tp")
		defer closer.Close()
		span := tracing.FromHttpGetSpan(tracer, context) // 生成span对象
		defer span.Finish()
		context.Set("tracer", tracer)
		context.Set("span", span)
		context.Next()
	}
}
