package router

import (
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"github/lunxun9527/bestpractice/pkg/xgin"
	"github/lunxun9527/bestpractice/server/accountApi/global"
	"github/lunxun9527/bestpractice/server/accountApi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap/zapcore"
)

func InitRouter(engine *gin.Engine) {

	engine.Use(gin.RecoveryWithWriter(zlog.NewWriter(zlog.DefaultLogger, zapcore.ErrorLevel), func(c *gin.Context, err any) {
		zlog.Errorf("recovery from panic err %v", err)
		xgin.FailWithLang(c)
	}))

	engine.Use(middleware.Trace(global.Config.JaegerTraceConf))
	engine.Use(func(c *gin.Context) {
		propagator := otel.GetTextMapPropagator()

		ctx := propagator.Extract(c, propagation.HeaderCarrier(c.Request.Header))

		zlog.InfofCtx(ctx, "request path %s ", c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	group := engine.Group("/api/v1")

	initAccountRouter(group)
}
