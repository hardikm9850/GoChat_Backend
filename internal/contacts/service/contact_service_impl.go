package service

import (
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	"github.com/hardikm9850/GoChat/internal/contacts/domain"
	"strings"
	dto "github.com/hardikm9850/GoChat/internal/contacts/handler/dto"
)

type ContactServiceImpl struct {
	userRepository repository.UserRepository
}

func (c *ContactServiceImpl) SyncContacts(userID string, request dto.SyncContactsRequest) ([]domain.ContactDTO, error) {
	phoneKeys := normalizeContacts(request.Contacts)

	if len(phoneKeys) == 0 {
		return []domain.ContactDTO{}, nil
	}

	users, err := c.userRepository.FindByMobiles(phoneKeys)
	if err != nil {
		return nil, err
	}

	contactDTOs := make([]domain.ContactDTO, 0, len(users))
	for _, user := range users {
		contactDTOs = append(contactDTOs, toContactDTO(user))
	}

	return contactDTOs, nil
}

func New(userRepository repository.UserRepository) ContactService {
	return &ContactServiceImpl{
		userRepository: userRepository,
	}
}

// Normalize and dedup contacts
func normalizeContacts(input []dto.ContactPayload) []domain.PhoneKey {
	set := make(map[string]domain.PhoneKey)

	for _, c := range input {
		phone := normalizePhone(c.Phone)
		code := strings.TrimSpace(c.CountryCode)

		if phone == "" || code == "" {
			continue
		}
		key := code + phone
		set[key] = domain.PhoneKey{
			CountryCode: code,
			Phone:       phone,
		}
	}
	res := make([]domain.PhoneKey, len(set))
	for _, v := range set {
		res = append(res, v)
	}
	return res
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	return phone
}
