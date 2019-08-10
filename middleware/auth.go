package middleware

import (
	"net/http"
	"net/url"

	echo "github.com/labstack/echo/v4"
)

type AuthConfig struct {
	signature func(url.Values, string) string
	secretKey string
}

func NewAuthConfig(signature func(url.Values, string) string, secretKey string) *AuthConfig {
	return &AuthConfig{
		signature: signature,
		secretKey: secretKey,
	}
}

var DefaultAuthConfig = &AuthConfig{}

func EchoAuth() echo.MiddlewareFunc {
	return EchoAuthWithConfig(DefaultAuthConfig)
}

// EchoAuth 用于 echo 框架的签名校验中间件
func EchoAuthWithConfig(authConfig *AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			formParams, err := ctx.FormParams()
			if err != nil {
				return ctx.String(http.StatusBadRequest, `400 Bad Request`)
			}
			sign := authConfig.signature(formParams, authConfig.secretKey)
			if sign != ctx.FormValue("sign") {
				return ctx.String(http.StatusBadRequest, `400 Bad Request`)
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
