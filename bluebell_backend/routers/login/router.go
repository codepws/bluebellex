// @Title  login
// @Description  登录模块路由
package login

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {

	/*
		// 需要签名验证，无需登录验证，无需 RBAC 权限验证
		login := e.Group("/api", r.middles.Signature())
		{
			login.POST("/login", signInHandler)
		}
	*/

	e.POST("/signin", signInHandler) //登录
	e.POST("/signup", signUpHandler) //注册
}
