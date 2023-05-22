package token

import (
	"fmt"
	"indoor_positioning/pkg/errno"
	"time"

	"github.com/zxmrlc/log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Context	JSON Web Token负载部分内容
type Context struct {
	ID uint64 // 用户ID
}

// @title	Sign
// @description	Sign signs the context with the specified secret.
// @auth	高宏宇
// @param	context	Context	自定义JWT负载	secret string JWT签名密钥
// @return	tokenString string JWT字符串	err error 错误信息
func Sign(context Context, secret string) (tokenString string, err error) {
	// 未指定签名密钥时，从配置文件载入默认密钥
	if secret == "" {
		secret = viper.GetString("jwt_secret")
	}

	//	The token content 官方字段
	//	iss (issuer)：签发人
	//	exp (expiration time)：过期时间
	//	sub (subject)：主题
	//	aud (audience)：受众
	//	nbf (Not Before)：生效时间
	//	iat (Issued At)：签发时间
	//	jti (JWT ID)：编号
	// 	token: 包含头部和负载的JWT对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  context.ID, // 自定义字段，用户ID
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Unix() + 7200, // 过期时间2h
	})

	// 使用密钥secret生成签名，生成最终令牌tokenString
	tokenString, err = token.SignedString([]byte(secret))

	return
}

// @title	ParseRequest
// @description	从参数中解析出token字符串，并进行token解析
// @auth	高宏宇
// @param	ctx *gin.Context
// @return	*Context JWT负载部分内容指针	error 错误信息
func ParseRequest(ctx *gin.Context) (*Context, error) {
	// 获取请求参数中的token字符串
	header := ctx.Request.Header.Get("Authorization")

	// 从配置文件读取JWT签名密钥
	secret := viper.GetString("jwt_secret")

	// 请求未添加token
	if len(header) == 0 {
		log.Error("ParseRequest error", errno.New(errno.ErrorTokenInvalid, fmt.Errorf("the length of the `Authorization` header is zero")))
		return &Context{}, errno.ErrorTokenInvalid
	}

	return Parse(header, secret)
}

// @title	Parse
// @description	根据给定密钥解析token字符串，如果token合法返回解析内容
// @auth	高宏宇
// @param	tokenString string token字符串	secret string JWT签名密钥
// @return	*Context JWT负载部分内容指针	error 错误信息
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	// 根据token字符串和签名密钥解析token，密钥来验证令牌的签名，确保令牌的完整性和真实性
	// 自定义secretFunc密钥提供函数按照实际需求提供密钥
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		return ctx, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = uint64(claims["id"].(float64))
		return ctx, nil

		// Other errors.
	} else {
		return ctx, err
	}
}

// @title	secretFunc
// @description	提供token签名验证密钥，同时检验token的加密方式
// @auth	高宏宇
// @param	secret string JWT签名密钥
// @return	jwt.Keyfunc JWT签名密钥提供函数
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// 类型断言，判断token.Method是否为*jwt.SigningMethodHMAC类型，即判断token使用的签名算法是否为HMAC算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}
