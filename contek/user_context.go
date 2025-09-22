package contek

import (
	"context"

	"github.com/SoeltanIT/agg-common-be/types"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserContext(ctx context.Context) *types.JWTClaims {
	user, ok := ctx.Value("user").(*jwt.Token)
	if !ok {
		return nil
	}

	return user.Claims.(*types.JWTClaims)
}

func GetUserRawToken(ctx context.Context) string {
	user := ctx.Value("user").(*jwt.Token)
	return user.Raw
}
