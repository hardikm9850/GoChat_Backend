package domain

type ContactDTO struct {
	ID          string `json:"user_id"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Name        string `json:"name"`
}

type SyncContactsResponse struct {
	RegisteredUsers []ContactDTO `json:"matched_users"`
	Timestamp       int64        `json:"timestamp"`
}
