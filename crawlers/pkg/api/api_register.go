package api

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"time"
)

func RegisterEndpoints() *gin.Engine {
	logger := zap.L()
	gin.SetMode(gin.ReleaseMode)
	var engine = gin.Default()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with local time format.
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, false))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	engine.Use(ginzap.RecoveryWithZap(logger, false))

	hd := NewHandler()
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.POST("/catalogs", hd.CreateCatalog)
	engine.POST("/sites", hd.CreateSite)
	engine.POST("/tasks/catalog-pages", hd.HandleCatalogPage)
	engine.POST("/tasks/novels", hd.HandleNovelPage)
	engine.POST("/tasks/schedule-task", hd.RunScheduleTask)

	return engine
}
