package service

import (
	authDomain "github.com/hardikm9850/GoChat/internal/auth/domain"
	contactsDomain "github.com/hardikm9850/GoChat/internal/contacts/domain"
)

func toContactDTO(user authDomain.User) contactsDomain.ContactDTO {
	return contactsDomain.ContactDTO{
		ID:          user.ID,
		Phone:       user.PhoneNumber,
		Name:        user.Name,
		CountryCode: user.CountryCode,
	}
}
