package core

import (
	"api-registry/init/config"
	"api-registry/init/db"
	"api-registry/init/log"
	"api-registry/init/redis"
	"api-registry/init/router"
)

func Init() {

	config.InitConfig()

	log.InitSysLog()

	redis.InitRedis()

	db.InitDB()

	router.InitRouter()
}
