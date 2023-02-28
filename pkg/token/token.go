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

// Context is the context of the JSON web token.
type Context struct {
	ID uint64
}

// Sign signs the context with the specified secret.
func Sign(context Context, secret string) (tokenString string, err error) {
	// Load the jwt secret from the Gin config if the secret isn't specified.
	if secret == "" {
		secret = viper.GetString("jwt_secret")
	}

	// The token content.
	/*	官方字段
		iss (issuer)：签发人
		exp (expiration time)：过期时间
		sub (subject)：主题
		aud (audience)：受众
		nbf (Not Before)：生效时间
		iat (Issued At)：签发时间
		jti (JWT ID)：编号
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  context.ID,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Unix() + 7200, // 过期时间2h
	})

	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(secret))

	return
}

// 从参数中解析出token字符串
// ParseRequest gets the token from the header and
// pass it to the Parse function to parses the token.
func ParseRequest(ctx *gin.Context) (*Context, error) {
	// 获取请求参数中的token字符串
	header := ctx.Request.Header.Get("Authorization")

	// Load the jwt secret from config
	secret := viper.GetString("jwt_secret")

	if len(header) == 0 {
		// TODO Errorf中如果以大写字母开头就会报“error strings should not be capitalized”，why
		log.Error("ParseRequest error", errno.New(errno.ErrorTokenInvalid, fmt.Errorf("the length of the `Authorization` header is zero")))
		return &Context{}, errno.ErrorTokenInvalid
	}

	return Parse(header, secret)
}

// Parse validates the token with the specified secret,
// and returns the context if the token was valid.
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	// Parse the token.
	// TODO 不太懂这里第2个参数的含义
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

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}
