package controller

import (
	"encoding/json"
	"github.com/swapbotAA/relayer/utils"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @Description //TODO 封装响应
 **/

type ResponseData struct {
	Code    utils.MyCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"` // omitempty当data为空时,不展示这个字段
}

func ResponseError(ctx *gin.Context, c utils.MyCode) {
	rd := &ResponseData{
		Code:    c,
		Message: c.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
	rd.printResponseData()
}

func ResponseErrorWithMsg(ctx *gin.Context, code utils.MyCode, msg string) {
	rd := &ResponseData{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
	rd.printResponseData()
}

func ResponseErrorWithData(ctx *gin.Context, code utils.MyCode, data interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
	rd.printResponseData()
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code:    utils.CodeSuccess,
		Message: utils.CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
	rd.printResponseData()
}

func (rd *ResponseData) printResponseData() {
	jsonData, _ := json.Marshal(rd)
	zap.L().Info("[Server response]", zap.String("response", string(jsonData)))
}
