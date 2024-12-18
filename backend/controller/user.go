package controller

import (
	"blog/service"
	"blog/utils"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	userController = newUser()
}

type user struct {
	secret     string
	payloadTtl time.Duration
}

func newUser() *user {
	return &user{
		secret:     viper.GetString("secret"),
		payloadTtl: time.Duration(viper.GetInt("payloadttlsec")) * time.Second,
	}
}

var userController *user

func GetUser() *user {
	return userController
}
func (u *user) LoginHandler(c *gin.Context) {
	var err error
	var verify = struct {
		Message   string `json:"message"`
		Signature string `json:"sign"`
	}{}
	err = c.BindJSON(&verify)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("参数出错"))
		return
	}
	if verify.Message == "" || verify.Signature == "" {
		c.JSON(200, utils.NewFailedResponse("登录失败"))
		return
	}
	var address string
	address, err = utils.Ecrecover(verify.Message, verify.Signature)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("签名出错"))
		return
	}
	claims := &utils.JwtCustomClaims{
		Address: address,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var auth string
	auth, err = token.SignedString([]byte(u.secret))
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("验证出错"))
		return
	}
	err = service.GetUser().AutoCreateIfNotExist(c, address, address)
	if err != nil {
		c.JSON(200, utils.NewFailedResponse("登陆失败"))
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponse(auth))
}
