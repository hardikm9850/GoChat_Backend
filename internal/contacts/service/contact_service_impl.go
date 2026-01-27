package service

import (
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	"github.com/hardikm9850/GoChat/internal/contacts/domain"
	dto "github.com/hardikm9850/GoChat/internal/contacts/handler/dto"
	internalDomain "github.com/hardikm9850/GoChat/internal/domain"
)

type ContactServiceImpl struct {
	userRepository repository.UserRepository
}

func (c *ContactServiceImpl) SyncContacts(userID string, request dto.SyncContactsRequest) (domain.SyncContactsResponse, error) {
	phoneKeys := internalDomain.NormalizeContacts(request.Contacts)

	if len(phoneKeys) == 0 {
		return domain.SyncContactsResponse{}, nil
	}

	users, err := c.userRepository.FindByMobiles(phoneKeys)
	if err != nil {
		return domain.SyncContactsResponse{}, err
	}

	contactDTOs := make([]domain.ContactDTO, 0, len(users))
	for _, user := range users {
		contactDTOs = append(contactDTOs, toContactDTO(user))
	}
	syncResponse := domain.SyncContactsResponse{
		RegisteredUsers: contactDTOs,
		Timestamp:       request.Timestamp,
	}
	return syncResponse, nil
}

func New(userRepository repository.UserRepository) ContactService {
	return &ContactServiceImpl{
		userRepository: userRepository,
	}
}
