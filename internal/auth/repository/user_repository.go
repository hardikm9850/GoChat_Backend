package repository

import (
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	contact_domain "github.com/hardikm9850/GoChat/internal/contacts/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	FindByID(id string) (domain.User, error)
	FindByMobile(mobile, countryCode string) (domain.User, error)
	FindByMobiles(mobile []contact_domain.PhoneKey) ([]domain.User, error)
}
