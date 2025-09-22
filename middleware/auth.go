package middleware

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/SoeltanIT/agg-common-be/types"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// NewAuthMiddlewareSignature : Initialize new instance for signature authentication middleware.
func NewAuthMiddlewareSignature(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		Claims:     &types.JWTClaimsSignature{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		},
	})
}

// NewAuthMiddleware : Initialize new instance for authentication middleware.
func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		Claims:     &types.JWTClaims{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		},
	})
}

func ValidateSocketToken(auth any) error {
	if auth != nil {
		result := make(map[string]string)
		v := reflect.ValueOf(auth)
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			result[key.String()] = fmt.Sprintf("%v", val.Interface())
		}

		tokenString, ok := result["token"]
		if ok {
			if tokenString == os.Getenv("STATIC_SECRET") {
				return nil
			}
		}
	}

	return errors.New("unauthorized")
}
