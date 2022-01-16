package shop

import (
	"bluebell_backend/common"
	"bluebell_backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

//商品模块查看商品路由Handler
func goodsHandler(c *gin.Context) {

	userID := c.GetUint64("Token_UserID")
	userVIP := c.GetInt("Token_UserVIP")
	var userClaims jwt.UserClaims
	userClaimsInterface, ok := c.Get("Token_User")
	if ok {
		userClaims = userClaimsInterface.(jwt.UserClaims)
	}

	common.ResponseSuccess(c, gin.H{
		"UserID":  userID,
		"UserVIP": userVIP,
		"User":    userClaims,
	})

}
