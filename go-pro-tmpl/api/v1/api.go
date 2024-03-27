package v1

import (
	"net/http"

	"api-registry/pkg/common"
	"api-registry/pkg/log"
	"api-registry/service"

	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context) {

	var apiService service.ApiService

	if err := ctx.ShouldBind(&apiService); err == nil {

		res := apiService.Hello(ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {

		log.SysLog.Errorf("Say Hello World get err: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
	}
}
