package router

import (
	v12 "api-registry/api/v1"
	"api-registry/middleware"

	"github.com/gin-gonic/gin"
)

func SettingRouter(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CorsMiddleware())

	v1 := r.Group("/api/v1")
	{
		// test ping request
		v1.GET("ping", v12.Hello)
	}

	return r
}
