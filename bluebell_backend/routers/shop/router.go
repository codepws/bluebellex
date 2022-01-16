// @Title  login
// @Description  登录模块路由
package shop

import (
	"bluebell_backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {

	// 认证路由组
	// authorized := r.Group("/", AuthRequired())
	// 和使用以下两行代码的效果完全一样:
	authorized := r.Group("/")
	// 路由组中间件! 在此例中，我们在 "authorized" 路由组中使用自定义创建的
	// AuthRequired() 中间件
	authorized.Use(jwt.JWTAuthMiddleware())
	{
		authorized.GET("/goods", goodsHandler) //商品

		// 嵌套路由组
		//testing := authorized.Group("testing")
		//testing.GET("/analytics", analyticsEndpoint)
	}

}
