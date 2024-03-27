package router

import (
	"api-registry/pkg/router"

	"github.com/gin-gonic/gin"
)

var ginRouter *gin.Engine

func InitRouter() {

	r := gin.Default()
	// router = gin.New()

	newRouter := router.SettingRouter(r)

	ginRouter = newRouter
}

func GetRouter() *gin.Engine {

	return ginRouter
}
