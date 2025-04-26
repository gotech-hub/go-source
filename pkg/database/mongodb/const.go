package mongodb

import (
	"os"
	"strings"
)

var (
	mappingRepositoryRegion = map[string]map[string][]string{
		"DEFAULT": {

			"config_active_games":     {"VN::vgr", "SEA::vgr"},
			"config_active_code":      {"VN::vgr", "SEA::vgr"},
			"config_active_reward":    {"VN::vgr", "SEA::vgr"},
			"redeem_code_transaction": {"VN::vgr", "SEA::vgr"},
			"config_vga_login":        {"VN::vgr", "SEA::vgr"},
			"code_management":         {"VN::vgr", "SEA::vgr"},
			"generation_request":      {"VN::vgr", "SEA::vgr"},

			"config_games":           {"VN::loyalty"},
			"config_currency":        {"VN::loyalty"},
			"lock_amount":            {"VN::loyalty"},
			"lock_amount_history":    {"VN::loyalty"},
			"order":                  {"VN::loyalty"},
			"point":                  {"VN::loyalty"},
			"point_transaction":      {"VN::loyalty"},
			"code_alias":             {"VN::loyalty"},
			"redeem_code_history":    {"VN::loyalty"},
			"privilege":              {"VN::loyalty"},
			"configs":                {"VN::loyalty"},
			"privilege_accumulation": {"VN::loyalty"},
			"privilege_benefit":      {"VN::loyalty"},
			"privilege_condition":    {"VN::loyalty"},
			"privilege_transaction":  {"VN::loyalty"},
			"tier":                   {"VN::loyalty"},
			"tier_version":           {"VN::loyalty"},
			"tier_history":           {"VN::loyalty"},
			"profile":                {"VN::loyalty"},
			"transaction_report":     {"VN::loyalty"},
			"schedule_report":        {"VN::loyalty"},
			"report":                 {"VN::loyalty"},
			"supported_game_product": {"VN::receiver"},
			"raw_payment_order":      {"VN::receiver"},
			"campaign":               {"VN::promotion"},
			"language":               {"VN::promotion"},
			"transaction_reward":     {"VN::promotion"},
			"compensation_history":   {"VN::promotion"},
			"redeem_history":         {"VN::promotion"},
			"reward_detail":          {"VN::promotion"},
			"reward":                 {"VN::promotion"},
			"stock":                  {"VN::promotion"},
			"query_template":         {"VN::promotion"},
			"template_config":        {"VN::promotion"},
			"template":               {"VN::promotion"},
			"rule":                   {"VN::promotion"},
			"product_metadata":       {"VN::promotion"},
			"product":                {"VN::promotion"},
			"code":                   {"VN::promotion"},
			"placement":              {"VN::promotion"},
			"promotion":              {"VN::promotion"},
		},
		"TEST": {

			"config_active_games":     {"VN::vgr_test", "SEA::vgr_test"},
			"config_active_code":      {"VN::vgr_test", "SEA::vgr_test"},
			"config_active_reward":    {"VN::vgr_test", "SEA::vgr_test"},
			"redeem_code_transaction": {"VN::vgr_test", "SEA::vgr_test"},
			"config_vga_login":        {"VN::vgr_test", "SEA::vgr_test"},
			"code_management":         {"VN::vgr", "SEA::vgr"},
			"generation_request":      {"VN::vgr", "SEA::vgr"},

			"config_games":           {"VN::loyalty_test"},
			"config_currency":        {"VN::loyalty_test"},
			"lock_amount":            {"VN::loyalty_test"},
			"lock_amount_history":    {"VN::loyalty_test"},
			"order":                  {"VN::loyalty_test"},
			"point":                  {"VN::loyalty_test"},
			"point_transaction":      {"VN::loyalty_test"},
			"code_alias":             {"VN::loyalty_test"},
			"redeem_code_history":    {"VN::loyalty_test"},
			"privilege":              {"VN::loyalty_test"},
			"configs":                {"VN::loyalty_test"},
			"privilege_accumulation": {"VN::loyalty_test"},
			"privilege_benefit":      {"VN::loyalty_test"},
			"privilege_condition":    {"VN::loyalty_test"},
			"privilege_transaction":  {"VN::loyalty_test"},
			"tier":                   {"VN::loyalty_test"},
			"tier_version":           {"VN::loyalty_test"},
			"tier_history":           {"VN::loyalty_test"},
			"profile":                {"VN::loyalty_test"},
			"transaction_report":     {"VN::loyalty_test"},
			"schedule_report":        {"VN::loyalty_test"},
			"report":                 {"VN::loyalty_test"},
			"supported_game_product": {"VN::receiver_test"},
			"raw_payment_order":      {"VN::receiver_test"},
			"campaign":               {"VN::promotion_test"},
			"language":               {"VN::promotion_test"},
			"transaction_reward":     {"VN::promotion_test"},
			"compensation_history":   {"VN::promotion_test"},
			"redeem_history":         {"VN::promotion_test"},
			"reward_detail":          {"VN::promotion_test"},
			"reward":                 {"VN::promotion_test"},
			"stock":                  {"VN::promotion_test"},
			"query_template":         {"VN::promotion_test"},
			"template_config":        {"VN::promotion_test"},
			"template":               {"VN::promotion_test"},
			"rule":                   {"VN::promotion_test"},
			"product_metadata":       {"VN::promotion_test"},
			"product":                {"VN::promotion_test"},
			"code":                   {"VN::promotion_test"},
			"placement":              {"VN::promotion_test"},
			"promotion":              {"VN::promotion_test"},
		},
	}
	countryMappingRegion = map[string]string{
		"VN":  "VN",
		"TW":  "SEA",
		"HK":  "SEA",
		"SG":  "SEA",
		"MY":  "SEA",
		"ID":  "SEA",
		"TH":  "SEA",
		"PH":  "SEA",
		"SEA": "SEA",
	}
)

func GetMappingRepositoryRegion(collectionName string) []string {
	env := strings.ToUpper(os.Getenv("ENV"))
	v, ok := mappingRepositoryRegion[env][collectionName]
	if !ok {
		return mappingRepositoryRegion["DEFAULT"][collectionName]
	}
	return v
}

func GetRegionCountry(country string) string {
	return countryMappingRegion[country]
}
