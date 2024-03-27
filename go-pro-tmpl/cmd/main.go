package main

import (
	"api-registry/init/router"
	"api-registry/pkg/core"
	"api-registry/pkg/log"
	"api-registry/pkg/model"

	"github.com/gin-gonic/gin"
)

func main() {

	core.Init()

	ginRouter := router.GetRouter()
	config := model.GetConfig()
	gin.SetMode(config.System.Application.Model)
	err := ginRouter.Run(config.System.Application.Port)
	if err != nil {
		log.SysLog.Warn("Api registry center is running failed, err: %v", err)
		return
	}

	log.SysLog.Info("Api registry center is running ...")

}
