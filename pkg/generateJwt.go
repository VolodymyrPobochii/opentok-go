package pkg

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJwt(config *Config /*config map[string]interface{}*/) (string, error) {
	now := time.Now().UnixNano() / int64(time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.ApiKey,
		"ist": "project",
		"iat": now,
		"exp": now + config.Auth.Expire,
	})
	stoken, err := token.SignedString([]byte(config.ApiSecret))
	if err != nil {
		return "", err
	}
	return stoken, nil
}
