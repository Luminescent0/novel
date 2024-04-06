package api

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"novel/dao"
	"novel/service"
	"novel/tool"
	"strings"
	"time"
)

const tokenSalt = "re_token"
const refreshExpiration = 60 * 60 * 24 * 7 * time.Second

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var MySecret = []byte("xian1")

func CreateToken(username string) (string, string, error) {
	//创建自己的声明
	aClaims := MyClaims{
		username,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*2, //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "novel",
			Subject:   "xian", //签发人
		},
	}
	//使用指定的签名方法创建对象
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims)
	token, err := reqClaim.SignedString(MySecret) //内部生成签名字符串，再用于获取完整、已签名的token
	if err != nil {
		return "", "", err
	}
	retoken := createRefreshToken(token, username)
	return token, retoken, nil
}
func createRefreshToken(accToken, username string) string {
	token := strings.Split(accToken, ".")
	payload := token[1]

	h := sha256.New()
	h.Write([]byte(payload + tokenSalt))
	retoken := fmt.Sprintf("%x", h.Sum(nil))

	go service.Set(retoken, username, refreshExpiration)
	err := dao.AddRefreshToken(username, retoken)
	if err != nil {
		log.Printf("set refresh token failed,username:%s,rt:%s,err:%s\n", username, retoken, err)
		retoken = ""
	}
	return retoken
}

//刷新token
func refreshToken(rt string) (acToken string, reToken string, err error) {
	username, err := service.Get(rt)
	if err != nil {
		log.Println("get rt from redis failed,err:", err)

		username = dao.GetRefreshToken(rt)
	}
	if username == "" {
		err = errors.New("refresh token failed")
		return
	}
	go service.Del(rt)
	go dao.DelRefreshToken(rt, username)

	return CreateToken(username)
}

func JwtAuthMiddleware(ctx *gin.Context) {
	//假设token放在Header的Authorization中，并使用Bearer开头
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		tool.RespSuccessfulWithData(ctx, gin.H{
			"msg": "请求头中auth为空",
		})
		ctx.Abort()
		return
	}
	parts := strings.SplitN(authHeader, "", 2)
	if len(parts) != 2 && parts[0] != "Bearer" {
		tool.RespSuccessfulWithData(ctx, gin.H{
			"msg": "请求头中auth格式有误",
		})
		ctx.Abort()
		return
	}
	var myClaims MyClaims
	token, err := jwt.ParseWithClaims(authHeader, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	fmt.Println(token)
	fmt.Println(token.Claims)
	var (
		claims *MyClaims
		ok     bool
	)

	if claims, ok = token.Claims.(*MyClaims); !ok {
		tool.RespSuccessfulWithData(ctx, gin.H{
			"msg": "token无效",
		})
		ctx.Abort()
		return
	}
	username := claims.Username
	fmt.Println(username)
	if err != nil {
		fmt.Println("parse token failed err", err)
		tool.RespInternalError(ctx)
		return
	}
	if token.Valid == false { //验证基于时间的声明
		tool.RespSuccessfulWithData(ctx, gin.H{
			"msg": "token过期",
		})
		ctx.Abort()
		return
	}
	//将当前请求的username信息保存到请求的上下文中
	ctx.Set("username", username)
	ctx.Next() //后续的处理函数可以通过ctx.Get()来获取当前请求的用户信息
}

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		if method != "" {
			ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("origin"))
			ctx.Header("Access-Control-Allow-Methods", "POST,GET,DELETE,OPTIONS,PUT")
			ctx.Header("Access-Control-Allow-Headers", "Origin,X-Requested-With,Content-Type,Accept,Authorization")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}
