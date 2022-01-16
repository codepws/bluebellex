package jwt

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"

	"bluebell_backend/common"
	"bluebell_backend/common/errno"
	"bluebell_backend/controller"
	"bluebell_backend/settings"
)

//用户鉴权
type UserClaims struct {
	UserID uint64 `json:"userid"`
	VIP    uint   `json:"vip"`
	jwt.StandardClaims
}

//Aeccse Token 过期时间
const APP_TOKEN_ACCESS_EXPIRE_DURATION = time.Minute * 10

//Refresh Token 过期时间
const APP_TOKEN_REFRESH_EXPIRE_DURATION = time.Hour * 24 * 7

var APP_TOKEN_SECRET_KEY = []byte("夏天夏天悄悄过去")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return APP_TOKEN_ACCESS_EXPIRE_DURATION, nil
}

// WrapToken 生成token，并封装JWT Token
func WrapToken(userID uint64, vip uint) (aToken string, rToken string, err error) {

	// 创建一个我们自己的声明
	c := UserClaims{
		userID, // 自定义字段
		vip,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(APP_TOKEN_ACCESS_EXPIRE_DURATION).Unix(), // 过期时间
			Issuer:    settings.WebConf.Name,                                   // 签发人(项目名称)
		},
	}

	// 加密并获得完整的编码后的字符串token
	// 使用指定的secret签名并获得完整的编码后的字符串token
	//Aeccse Token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(APP_TOKEN_SECRET_KEY)
	if err != nil {
		return
	}
	fmt.Println(aToken)

	//Refresh Token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(APP_TOKEN_REFRESH_EXPIRE_DURATION).Unix(),
		Issuer:    settings.WebConf.Name,
	}).SignedString(APP_TOKEN_SECRET_KEY)
	if err != nil {
		return
	}

	return
}

// 解析token
func UnwrapToken(tokenStr string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return APP_TOKEN_SECRET_KEY, nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	/*
		var claims = new(UserClaims)
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return APP_TOKEN_SECRET_KEY, nil
		})

		if token != nil { // 校验token
			if token.Valid {
				return claims, nil
			} else {
				//return nil, errors.New("invalid token")
			}
		}
	*/

	fmt.Println("jwt.ParseWithClaims error:", err) // token is expired by 1h16m29s
	return nil, err
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authTokenValue := c.Request.Header.Get("X-Token")
		if authTokenValue == "" {
			common.ResponseErrorWithMsg(c, -1, "请求头缺少Auth Token", nil)
			c.Abort()
			return
		}

		fmt.Println(authTokenValue)
		// 按空格分割
		//parts := strings.SplitN(authHeader, " ", 2)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		//	c.Abort()
		//	return
		//}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := UnwrapToken(authTokenValue)
		if err != nil {
			fmt.Println(err)
			common.ResponseError(c, errno.ErrTokenValidation)
			c.Abort() //用来终止request, 阻止其到达handler 一般情况下用在鉴权与认证的中间件中
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(controller.KeyCtxContextUserID, mc.UserID)
		//c.Set(controller.KeyCtxUserName, mc.UserName)
		c.Set(controller.KeyCtxContextUserVIP, mc.VIP)
		c.Set(controller.KeyCtxContextUser, *mc)

		// 后续的处理函数可以用过c.Get("Token_UserID")来获取当前请求的用户信息
		c.Next()
	}
}

/*
func AuthMiddleware(next http.Handler) http.Handler {
    if len(APP_KEY) == 0 {
        log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
    }
    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(APP_KEY), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })
    return jwtMiddleware.Handler(next)
}
*/
