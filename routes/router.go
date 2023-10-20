package routes

import (
	"bluebellproject/logger"
	"bluebellproject/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", setting.Conf.AppConfig.Version)
	})
	return r
}
