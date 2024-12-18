package jwt

import (
	"blog/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

// address是只增不减
func NewVerifyMiddleware() func(*gin.Context) {
	return func(ctx *gin.Context) {
		rawauthorization := ctx.GetHeader("Authorization")
		authorization, ok := strings.CutPrefix(rawauthorization, "Bearer ")
		if !ok {
			ctx.AbortWithStatusJSON(401, utils.NewFailedResponse("未登录"))
			return
		}
		claims := utils.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(authorization, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("secret")), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(401, utils.NewFailedResponse("未登录"))
			return
		}
		ctx.Set("address", claims.Address)
		ctx.Next()
	}
}
