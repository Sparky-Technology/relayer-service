package routers

import (
	"github.com/swapbotAA/relayer/controller"
	_ "github.com/swapbotAA/relayer/docs"
	"github.com/swapbotAA/relayer/logger"
	"github.com/swapbotAA/relayer/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-contrib/pprof"
)

// SetupRouter 设置路由
func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}
	//初始化 gin Engine  新建一个没有任何默认中间件的路由
	r := gin.New()
	//设置中间件
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true),                           // Recovery 中间件会 recover掉项目可能出现的panic，并使用zap记录相关日志
		middlewares.RateLimitMiddleware(2*time.Second, 40), // 每两秒钟添加十个令牌  全局限流
		middlewares.Cors(),
	)

	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// 注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	// 注册登陆业务
	v1.POST("/login", controller.LoginHandler)               // 登陆业务
	v1.POST("/signup", controller.SignUpHandler)             // 注册业务
	v1.GET("/refresh_token", controller.RefreshTokenHandler) // 刷新accessToken
	// swapbotAA
	v1.POST("/handle_op", controller.HandleOp)     // 快速交易
	v1.POST("/limit_order", controller.LimitOrder) // 限价单
	// 中间件
	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	pprof.Register(r) // 注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
