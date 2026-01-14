package http

type ContactSyncResponse struct {
	MatchedUsers []MatchedUser `json:"matched_users"`
}

type MatchedUser struct {
	UserID      string `json:"user_id"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Username    string `json:"username"`
}
