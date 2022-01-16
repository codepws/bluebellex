package login

import (
	"fmt"
	"net/http"

	"bluebell_backend/common"
	"bluebell_backend/common/errno"
	"bluebell_backend/controller"
	"bluebell_backend/models/login"
	"bluebell_backend/pkg/jwt"
	"bluebell_backend/pkg/logger"
	service_login "bluebell_backend/service/login"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	//"gopkg.in/go-playground/validator.v8"
	//"github.com/go-playground/validator/v10"
	//"github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
)

//登录模块登录路由Handler
func signInHandler(c *gin.Context) {

	// 1.获取请求参数
	// 2.校验数据有效性
	var signInRequest login.SignInRequest
	//ShouldBindJSON
	if err := c.ShouldBind(&signInRequest); err != nil {
		//返回响应
		common.ResponseError(c, errno.ErrParam)
		return
	}

	//var user login.User
	// 参数接收正确, 进行登录操作
	if user, err := service_login.Login(&signInRequest); err != nil {
		logger.MainLogger.Error("service.Login(&u) failed", zap.Error(err))
		common.ResponseError(c, errno.ErrNameOrPasswordIncorrect)
		return
	} else {
		if aToken, rToken, err := jwt.WrapToken(user.UserID, user.VIP); err != nil {

			// 比如，返回“用户手机号不合法”错误
			//c.JSON(http.StatusOK, errno.ErrTokenSign.WithID(c.GetString("trace-id")))

			//返回响应
			common.ResponseError(c, errno.ErrTokenSign.WithData("生成Token失败！"))

		} else {

			//在接收request请求时就设置好response时要添加的header
			c.Header(controller.KeyCtxHeaderToken, aToken)
			//c.Header("X-RToken", rToken)

			//返回响应
			common.ResponseSuccess(c, gin.H{
				//"accessToken":  aToken,
				"refreshToken": rToken,
				"user":         user,
			})

			//c.JSON(http.StatusOK, errno.ErrTokenSign.WithData("登录成功").WithID(c.GetString("trace-id")))
		}
	}

}

//登录模块注册路由Handler
func signUpHandler(c *gin.Context) {

	// 1.获取请求参数
	// 2.校验数据有效性
	var signUpRequest login.SignUpRequest
	//ShouldBindJSON
	if err := c.ShouldBind(&signUpRequest); err != nil {

		fmt.Printf("=================>Error:  %T\n", err)

		logger.MainLogger.Error("signUpHandler.Bind(&u) failed", zap.Error(err))
		// 验证错误
		/*
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error() })
		*/

		// 验证错误
		//c.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"message": signUpRequest.GetError(err.(validator.ValidationErrors)), // 注意这里要将 err 进行转换
		//})

		//把类型判断的结果分配给变量，这个变量就变成了类型转换后的结果，这样在对应case里处理的时候，就不用再来一次类型转换了
		switch err := err.(type) {
		case validator.ValidationErrors:
			common.ResponseErrorWithMsg(c, http.StatusUnprocessableEntity, "", gin.H{
				"message": signUpRequest.GetError(err), // 注意这里要将 err 进行转换  err.(validator.ValidationErrors)
			})
		default:
			// {"error": "strconv.ParseUint: parsing \"aaaaaa\": invalid syntax"}

			common.ResponseError(c, errno.ErrParam)
		}

		return
	}

	// 参数接收正确, 进行注册操作
	// 3.注册用户
	//var user login.User
	// 参数接收正确, 进行登录操作
	if user, err := service_login.Register(&signUpRequest); err != nil {
		logger.MainLogger.Error("service.Login(&u) failed", zap.Error(err))
		common.ResponseError(c, errno.ErrNameOrPasswordIncorrect)
		return
	} else {

		if aToken, rToken, err := jwt.WrapToken(user.UserID, user.VIP); err != nil {

			logger.MainLogger.Error("signUpHandler.WrapToken(&u) failed", zap.Error(err))

			//返回响应
			common.ResponseError(c, errno.ErrTokenSign.WithData("生成Token失败！"))

		} else {

			//在接收request请求时就设置好response时要添加的header
			c.Header(controller.KeyCtxHeaderToken, aToken)

			//返回响应
			common.ResponseSuccess(c, gin.H{
				"refreshToken": rToken,
				"user":         user,
			})

		}
	}

}
