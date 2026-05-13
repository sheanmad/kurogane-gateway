package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kurogane/gateway/internal/handler"
	"github.com/kurogane/gateway/internal/middleware"
	"github.com/kurogane/gateway/internal/proxy"
)

func Setup(p *proxy.EngineProxy) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Auth())
	r.Use(middleware.RateLimit())
	r.Use(middleware.Errors())

	runHandler := handler.NewRunHandler(p)

	r.GET("/health", handler.HealthCheck)
	r.POST("/runs", runHandler.CreateRun)
	r.GET("/runs", runHandler.ListRuns)
	r.GET("/runs/:id", runHandler.GetRun)
	r.GET("/runs/:id/events", runHandler.StreamRunEvents)

	return r
}