package controllers

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"ottodigital.id/library/ottotracing"
)

func TracingFirstControllerCtx(c *gin.Context, request interface{}, namectrl string) opentracing.Span {
	var span opentracing.Span
	if cspan, ok := c.Get("tracing-context"); ok {
		span = ottotracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), namectrl, c.Request.Method, c.Request.URL.Path)

	} else {
		span = ottotracing.StartSpanWithHeader(&c.Request.Header, c.Request.Method, namectrl, c.Request.URL.Path)
	}

	data, _ := json.Marshal(request)
	if len(data) > constants.MAXUDP {
		request = fmt.Sprint("%s", data[:constants.MAXUDP])
	}

	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("REQ", request),
		zap.Duration("backoff", time.Second))
	return span
}

func TracingEmptyFirstControllerCtx(c *gin.Context, namectrl string) opentracing.Span {
	var span opentracing.Span
	if cspan, ok := c.Get("tracing-context"); ok {
		span = ottotracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), namectrl, c.Request.Method, c.Request.URL.Path)

	} else {
		span = ottotracing.StartSpanWithHeader(&c.Request.Header, c.Request.Method, namectrl, c.Request.URL.Path)
	}
	return span
}
