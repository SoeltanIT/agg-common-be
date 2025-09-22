package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID              string         `json:"id"`
	Namespace       string         `json:"namespace,omitempty"`
	ParentNamespace string         `json:"parent_namespace,omitempty"`
	UserId          string         `json:"user_id,omitempty"` // this is for when agent login value for super agent id
	Email           string         `json:"email"`
	Type            Role           `json:"position_type"`
	Permissions     PermissionsDTO `json:"permissions"`
	jwt.RegisteredClaims
}

type JWTClaimsSignature struct {
	ClientId     string `json:"client_id"`
	UserId       string `json:"user_id"`
	SuperAgentId string `json:"super_agent_id"`
	ClientSecret string `json:"client_secret"`
	jwt.RegisteredClaims
}
