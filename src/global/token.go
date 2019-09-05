package global

import (
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

type TokenClaims struct {
	UserID int64
	jwt.Payload
}

func GenToken(id int64) (string, error) {
	var tc TokenClaims
	tc.UserID = id

	// 到期时间点
	tc.ExpirationTime = jwt.NumericDate(time.Now().Add(time.Duration(Config.Token.Expire) * time.Hour))

	// 生成算法
	alg := jwt.NewHS256([]byte(Config.Token.Key))

	// 生成签名
	tokenStr, err := jwt.Sign(tc, alg)
	if err != nil {
		return "", err
	}

	return string(tokenStr), nil
}

func ParseToken(tokenStr string) (*TokenClaims, error) {
	var tc TokenClaims

	// 生成算法
	alg := jwt.NewHS256([]byte(Config.Token.Key))

	// 验证签名
	_, err := jwt.Verify([]byte(tokenStr), alg, &tc,
		// 验证过期时间
		jwt.ValidatePayload(&tc.Payload, jwt.ExpirationTimeValidator(time.Now())))
	if err != nil {
		return nil, err
	}

	return &tc, nil
}
