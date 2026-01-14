package http

// SyncContactsRequest Exported request struct for Swaggo
type SyncContactsRequest struct {
    Contacts  []ContactPayload `json:"contacts" binding:"required,min=1"`
    Timestamp int64            `json:"timestap"`
}

type ContactPayload struct {
    Phone       string `json:"phone" binding:"required,e164"`
    CountryCode string `json:"country_code" binding:"required"`
    Name        string `json:"name"`
}
