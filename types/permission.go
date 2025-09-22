package types

type PermissionsDTO struct {
	Dashboard                  int `json:"dashboard"`
	ReportPlayerActive         int `json:"report_player_active"`
	ReportClients              int `json:"report_clients"`
	ReportSlot                 int `json:"report_slot"`
	ReportProfit               int `json:"report_profit"`
	ReportClientShared         int `json:"report_client_shared"`
	SuperAgent                 int `json:"super_agent"`
	Agent                      int `json:"agent"`
	GameProviders              int `json:"game_providers"`
	Games                      int `json:"games"`
	PlayerPendingTransaction   int `json:"player_pending_transaction"`
	Settings                   int `json:"settings"`
	PermissionRegenerateSecret int `json:"permission_regenerate_secret"`
}
