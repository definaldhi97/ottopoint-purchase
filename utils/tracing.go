package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"ottodigital.id/library/ottotracing"
)

func TracingCtx(c *gin.Context, namectrl string) opentracing.Span {
	var span opentracing.Span
	if cspan, ok := c.Get("tracing-context"); ok {
		span = ottotracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), namectrl, c.Request.Method, c.Request.URL.Path)

	} else {
		fmt.Println("Woy")
		span = ottotracing.StartSpanWithHeader(&c.Request.Header, c.Request.Method, namectrl, c.Request.URL.Path)
	}
	zaplog.InfoWithSpan(span, "Test Span",
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
	return span
}
