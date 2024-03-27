package service

import (
	"api-registry/pkg/log"
	"context"

	"api-registry/pkg/common"
	"api-registry/pkg/message"
	"api-registry/pkg/model"
)

type ApiService struct {
	ApiService model.APIRegistryInfo
}

func (as ApiService) Hello(ctx context.Context) common.Response {

	var code = message.Success

	log.SysLog.Info("The System say: 'Hello World!'")

	return common.Response{
		Status: code,
		Msg:    message.GetMsg(code),
		Data:   "This is a test, Hello World!",
	}
}
