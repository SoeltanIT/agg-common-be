package types

import (
	"net/http"
)

type Permission string

// Role represents the role of a user
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleClient Role = "client"
	RoleUser   Role = "user"
)
const (
	PermissionDashboard          Permission = "dashboard"
	PermissionReportPlayerActive Permission = "report_player_active"
	PermissionReportClients      Permission = "report_clients"
	PermissionReportSlot         Permission = "report_slot"
	PermissionReportProfit       Permission = "report_profit"
	PermissionReportClientShared Permission = "report_client_shared"
	PermissionSuperAgent         Permission = "super_agent"
	PermissionAgent              Permission = "agent"
	PermissionGameProviders      Permission = "game_providers"
	PermissionGames              Permission = "games"
	PermissionPlayerPendingTxn   Permission = "player_pending_transaction"
	PermissionSettings           Permission = "settings"
	PermissionRegenerateSecret   Permission = "permission_regenerate_secret"
)

type PermissionAction int

const (
	ActionRead   PermissionAction = 1 << iota // 1
	ActionWrite                               // 2
	ActionDelete                              // 4
)

var MethodAction = map[string]PermissionAction{
	http.MethodGet:    ActionRead,
	http.MethodPost:   ActionWrite,
	http.MethodPut:    ActionWrite,
	http.MethodPatch:  ActionWrite,
	http.MethodDelete: ActionDelete,
}
