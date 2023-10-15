package controller

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/swapbotAA/relayer/dao/mysql"
	"github.com/swapbotAA/relayer/models"
	"github.com/swapbotAA/relayer/pkg/jwt"
	"github.com/swapbotAA/relayer/service"
	"github.com/swapbotAA/relayer/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler
// @Summary 注册业务
// @Description 注册业务
// @Tags 用户业务接口
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.RegisterForm true "注册请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData
// @Router /signup [POST]
func SignUpHandler(c *gin.Context) {
	// 1.获取请求参数
	var fo *models.RegisterForm
	// 2.校验数据有效性
	if err := c.ShouldBindJSON(&fo); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithData(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return // 翻译错误
	}
	fmt.Printf("fo: %v\n", fo)
	// 3.业务处理 —— 注册用户
	if err := service.SignUp(fo); err != nil {
		zap.L().Error("service.signup failed", zap.Error(err))
		if err.Error() == mysql.ErrorUserExit {
			ResponseError(c, utils.CodeUserExist)
			return
		}
		ResponseError(c, utils.CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler
// @Summary 登录业务
// @Tags 用户业务接口
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.LoginForm true "登录请求参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData
// @Router /login [POST]
func LoginHandler(c *gin.Context) {
	// 1、获取请求参数及参数校验
	var u *models.LoginForm
	if err := c.ShouldBindJSON(&u); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithData(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2、业务逻辑处理——登录
	user, err := service.Login(u)
	if err != nil {
		zap.L().Error("service.Login failed", zap.String("username", u.UserName), zap.Error(err))
		if err.Error() == mysql.ErrorUserNotExit {
			ResponseError(c, utils.CodeUserNotExist)
			return
		}
		ResponseError(c, utils.CodeInvalidParams)
		return
	}
	// 3、返回响应
	ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID), //js识别的最大值：id值大于1<<53-1  int64: i<<63-1
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

// RefreshTokenHandler
// @Summary 刷新accessToken
// @Description 刷新accessToken
// @Tags 用户业务接口
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseData
// @Router /refresh_token [GET]
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的 Authorization 中，并使用 Bearer 开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithData(c, utils.CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithData(c, utils.CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
