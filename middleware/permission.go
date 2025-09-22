package middleware

import (
	"encoding/json"
	"log/slog"
	"strings"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/SoeltanIT/agg-common-be/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getPermValue(perms types.PermissionsDTO, p types.Permission) int {
	switch p {
	case types.PermissionDashboard:
		return perms.Dashboard
	case types.PermissionReportPlayerActive:
		return perms.ReportPlayerActive
	case types.PermissionReportClients:
		return perms.ReportClients
	case types.PermissionReportSlot:
		return perms.ReportSlot
	case types.PermissionReportProfit:
		return perms.ReportProfit
	case types.PermissionReportClientShared:
		return perms.ReportClientShared
	case types.PermissionSuperAgent:
		return perms.SuperAgent
	case types.PermissionAgent:
		return perms.Agent
	case types.PermissionGameProviders:
		return perms.GameProviders
	case types.PermissionGames:
		return perms.Games
	case types.PermissionPlayerPendingTxn:
		return perms.PlayerPendingTransaction
	case types.PermissionSettings:
		return perms.Settings
	case types.PermissionRegenerateSecret:
		return perms.PermissionRegenerateSecret
	default:
		return 0
	}
}

func parsePermissionFromQuery(q string) (types.Permission, bool) {
	switch strings.ToLower(strings.TrimSpace(q)) {
	case "dashboard":
		return types.PermissionDashboard, true
	case "report_player_active":
		return types.PermissionReportPlayerActive, true
	case "report_clients":
		return types.PermissionReportClients, true
	case "report_slot":
		return types.PermissionReportSlot, true
	case "report_profit":
		return types.PermissionReportProfit, true
	case "report_client_shared":
		return types.PermissionReportClientShared, true
	case "super_agent":
		return types.PermissionSuperAgent, true
	case "agent":
		return types.PermissionAgent, true
	case "game_providers":
		return types.PermissionGameProviders, true
	case "games":
		return types.PermissionGames, true
	case "player_pending_transaction":
		return types.PermissionPlayerPendingTxn, true
	case "settings":
		return types.PermissionSettings, true
	case "permission_regenerate_secret":
		return types.PermissionRegenerateSecret, true
	default:
		return "", false
	}
}

func ValidatePermission(requiredPermissions ...types.Permission) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		u := c.Locals("user")
		if u == nil {
			slog.Warn("ValidatePermission: JWT not found in Locals")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		token, ok := u.(*jwt.Token)
		if !ok || token == nil {
			slog.Warn("ValidatePermission: JWT token invalid type")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok || claims == nil {
			slog.Warn("ValidatePermission: JWT claims invalid or nil")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		if b, _ := json.Marshal(claims.Permissions); len(b) > 0 {
			slog.Debug("ValidatePermission: permissions snapshot", "permissions", string(b))
		}

		method := c.Method()
		action, ok := types.MethodAction[method]
		if !ok {
			slog.Warn("ValidatePermission: method not mapped to action", "method", method)
			return common.Response().SetError(common.ErrForbidden).Send(c)
		}

		queryRole := c.Query("role")
		queryType := c.Query("type")
		query := queryRole
		if query == "" {
			query = queryType
		}

		if permFromQuery, has := parsePermissionFromQuery(query); has {
			if len(requiredPermissions) > 0 {
				foundInRequired := false
				for _, rp := range requiredPermissions {
					if rp == permFromQuery {
						foundInRequired = true
						break
					}
				}
				if !foundInRequired {
					slog.Info("Permission denied - queried permission not in required list",
						"userID", claims.ID, "query", query, "method", method)
					return common.Response().SetError(common.ErrForbidden).Send(c)
				}
			}

			permVal := getPermValue(claims.Permissions, permFromQuery)
			if (permVal & int(action)) != 0 {
				slog.Info("Permission granted",
					"userID", claims.ID, "permission", string(permFromQuery),
					"action", action, "method", method, "query", query)
				return c.Next()
			}

			slog.Info("Permission denied",
				"userID", claims.ID, "permission", string(permFromQuery),
				"action", action, "method", method, "query", query, "value", permVal)
			return common.Response().SetError(common.ErrForbidden).Send(c)
		}

		if len(requiredPermissions) == 0 {
			return c.Next()
		}

		for _, p := range requiredPermissions {
			permVal := getPermValue(claims.Permissions, p)
			if (permVal & int(action)) != 0 {
				slog.Info("Permission granted",
					"userID", claims.ID, "permission", string(p),
					"action", action, "method", method)
				return c.Next()
			}
			slog.Debug("Permission not sufficient",
				"userID", claims.ID, "permission", string(p),
				"have", permVal, "needAction", action, "method", method)
		}

		slog.Info("Permission denied - no valid permissions found",
			"userID", claims.ID, "method", method, "required", requiredPermissions)
		return common.Response().SetError(common.ErrForbidden).Send(c)
	}
}

func ValidatePermissionUserClient() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			slog.Warn("JWT not found in Locals")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			slog.Warn("JWT Token invalid type")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok || claims == nil {
			slog.Warn("JWT Claims invalid or nil")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		return c.Next()
	}
}
func ValidateAggregatorSignature() fiber.Handler {
	return func(c *fiber.Ctx) error {
		signature := c.Get("X-Aggregator-Signature")
		if signature == "" {
			slog.Warn("X-Aggregator-Signature header not found")
			return common.Response().SetError(common.ErrMissingAggregatorSignature).Send(c)
		}

		user := c.Locals("user")
		if user == nil {
			slog.Warn("JWT not found in Locals")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			slog.Warn("JWT Token invalid type")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		claims, ok := token.Claims.(*types.JWTClaimsSignature)
		if !ok || claims == nil {
			slog.Warn("JWT Claims invalid or nil")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		clientID := claims.ClientId // keep this consistent
		if clientID == "" {
			slog.Warn("Client ID not found in JWT claims")
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}
		slog.Info("Validating signature for client_id", "client_id", clientID)

		/* TODO: access integrator agent from valkey
		agentRepository := repository.NewIntegratorDao(database.DB)
		agent, err := agentRepository.GetById(c.Context(), clientID)
		if err != nil {
			slog.Warn("Agent not found", "client_id", clientID, "error", err.Error())
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}

		if agent.User.Status != "active" {
			slog.Warn("Agent is not active", "client_id", clientID, "status", agent.User.Status)
			return common.Response().SetError(common.ErrUnauthorized).Send(c)
		}
		*/

		if err := ValidateSignature(signature); err != nil {
			slog.Warn("Signature validation failed", "client_id", clientID, "error", err.Error())

			if strings.Contains(err.Error(), "signature expired") {
				return common.Response().SetError(common.ErrSessionExpired).Send(c)
			}

			return common.Response().SetError(common.ErrInvalidSignature).Send(c)
		}

		c.Locals("user", claims)

		//slog.Info("Signature validation successful", "client_id", agent.ID)
		return c.Next()
	}
}
