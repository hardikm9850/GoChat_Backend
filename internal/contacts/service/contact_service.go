package service

import (
	"github.com/hardikm9850/GoChat/internal/contacts/domain"
	dto "github.com/hardikm9850/GoChat/internal/contacts/handler/dto"
)

type ContactService interface {
	// SyncContacts takes a list of phone numbers from the client
	// and returns registered users as DTOs
	SyncContacts(userID string, request dto.SyncContactsRequest) (domain.SyncContactsResponse, error)
}
