package middleware

import (
	"github.com/gin-gonic/gin"
	"server/service"
	"server/utils"
)

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"strconv"
)

var jwtService = service.ServiceGroupApp.JwtService

// 我来看看jwt验证是怎么回事？
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := utils.GetAccessToken(c)   // 从请求头中获取accessToken
		refreshToken := utils.GetRefreshToken(c) //  从cookic中获取refreshToken

		if jwtService.IsInBlacklist(refreshToken) { // 这个token是否在黑名单中
			utils.ClearRefreshToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()
			return
		}

		j := utils.NewJWT()

		claims, err := j.ParseAccessToken(accessToken) // 解析accessToken
		if err != nil {
			if accessToken == "" || errors.Is(err, utils.TokenExpired) {
				refreshClaims, err := j.ParseRefreshToken(refreshToken) // 解析refreshToken，获取用户信息，看看能不能用来获取一个新的accesstoken
				if err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("Refresh token expired or invalid", c)
					c.Abort()
					return
				}

				var user database.User
				if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("The user does not exist", c)
					c.Abort()
					return
				}

				newAccessClaims := j.CreateAccessClaims(request.BaseClaims{
					UserID: refreshClaims.UserID,
					UUID:   user.UUID,
					RoleID: user.RoleID,
				})

				newAccessToken, err := j.CreateAccessToken(newAccessClaims)
				if err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("Faild to create new access token", c)
					c.Abort()
					return
				}

				c.Header("new-access-token", newAccessToken)
				c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))

				c.Set("claims", &newAccessClaims) // 存储在gin的上下文中允许后续的使用
				c.Next()
				return
			}
			utils.ClearRefreshToken(c)
			response.NoAuth("Invalid access token", c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
