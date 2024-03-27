package log

import (
	"api-registry/pkg/log"
	"api-registry/pkg/model"
)

var config model.Config

func InitSysLog() {

	config = model.GetConfig()

	log.InitLogger(
		config.System.Application.LogFilePath,
		config.System.Application.LogLevel,
		config.System.Application.ConsoleLog)

	log.SysLog.Info("Init SysLog Success!")

}
