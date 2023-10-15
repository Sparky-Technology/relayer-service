package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swapbotAA/relayer/service/limitorder/manager"
	"github.com/swapbotAA/relayer/utils"
	"go.uber.org/zap"

	"github.com/swapbotAA/relayer/service/limitorder/bo"
)

// LimitOrder
// @Title Limit Order
// @Description limit order operation
// @Tags swapbotAA业务接口
// @Param object body vo.LimitOrderVO true "限价单参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData
// @Router /limit_order [POST]
func LimitOrder(c *gin.Context) {
	// 处理请求参数
	var bo *bo.LimitOrderBO
	if err := c.ShouldBindJSON(&bo); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("[LimitOrder] Invalid param!", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithData(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans))) // 翻译错误
		return
	}
	zap.L().Info("[LimitOrder] BO:", zap.String("bo", bo.ToString()))
	// 业务处理
	txHash, err := manager.AddOrder(bo)
	if err != nil {
		zap.L().Error("[LimitOrder] Execution reverted!", zap.Error(err))
		if myError, ok := err.(*utils.MyError); ok {
			ResponseErrorWithData(c, utils.MyCode(myError.Code), myError.Msg)
			return
		}
		ResponseError(c, utils.CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, txHash)
}
