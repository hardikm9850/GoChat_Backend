package domain

import "time"

type User struct {
	ID           string
	PhoneNumber  string
	CountryCode  string
	Name         string
	PasswordHash string
	PhoneHash    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
