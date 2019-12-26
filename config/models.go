package config

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type OttoConfig struct {
	Ctx *gin.Context
	Idlog string
	Parentspan opentracing.Span
	Childspan opentracing.Span
	Request interface{}
	Response interface{}
}
