package memory

import (
	"errors"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	domain2 "github.com/hardikm9850/GoChat/internal/contacts/domain"
	"sync"
)

type UserRepo struct {
	mu       sync.RWMutex
	users    map[string]domain.User // key: userID
	byMobile map[string]string      // mobile -> userID
}

func New() *UserRepo {
	return &UserRepo{
		users:    make(map[string]domain.User),
		byMobile: make(map[string]string),
	}
}

func (r *UserRepo) Create(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byMobile[user.PhoneNumber]; exists {
		return errors.New("user already exists in memory")
	}

	r.users[user.ID] = user
	r.byMobile[user.PhoneNumber] = user.ID
	return nil
}

func (r *UserRepo) FindByMobile(mobile, countryCode string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byMobile[mobile]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	user := r.users[id]
	return user, nil
}

func (r *UserRepo) FindByMobiles(mobile []domain2.PhoneKey) ([]domain.User, error) {
	return nil, nil
}

func (r *UserRepo) FindByID(id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil
}
