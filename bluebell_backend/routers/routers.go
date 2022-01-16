package routers

import (
	"bluebell_backend/routers/login" //同一项目下
	"bluebell_backend/routers/shop"  //同一项目下

	"github.com/gin-gonic/gin"
)

type RouterOption func(*gin.Engine)

var routerOptions = []RouterOption{}

// 设计模式–函数式选项模式
// 注册app的路由配置
func include(opts ...RouterOption) {
	routerOptions = append(routerOptions, opts...)
}

//初始化
func SetupRouter() *gin.Engine {

	////////////////////////////////////////
	// 加载多个APP的路由配置
	include(login.Routers, shop.Routers)

	////////////////////////////////////////
	r := gin.Default()
	//r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//1.首位多余元素会被删除(../ or //);
	//2.然后路由会对新的路径进行不区分大小写的查找;
	//3.如果能正常找到对应的handler，路由就会重定向到正确的handler上并返回301或者307.(比如: 用户访问/FOO 和 /..//Foo可能会被重定向到/foo这个路由上)
	r.RedirectFixedPath = true
	for _, opt := range routerOptions {
		opt(r)
	}

	return r
}
