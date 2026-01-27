package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hardikm9850/GoChat/internal/contacts/domain"
	dto "github.com/hardikm9850/GoChat/internal/contacts/handler/dto"
	"regexp"
	"strings"
)

// HashPhoneNumber creates a deterministic hash for contact matching
func HashPhoneNumber(countryCode, phone string) string {
	normalized := countryCode + strings.ReplaceAll(phone, " ", "")
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// NormalizePhoneNumber ensures a consistent format
func NormalizePhoneNumber(countryCode, phoneNumber string) string {
	// Remove any spaces, dashes, parentheses
	cleaned := regexp.MustCompile(`[^0-9+]`).ReplaceAllString(phoneNumber, "")
	return countryCode + cleaned
}

// Normalize and de-dup contacts
func NormalizeContacts(input []dto.ContactPayload) []domain.PhoneKey {
	set := make(map[string]domain.PhoneKey)

	for _, c := range input {
		phone := normalizePhone(c.Phone)
		//code := strings.TrimSpace(c.CountryCode)
		code := c.CountryCode

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

func NormalizeCountryCode(input string) (string, error) {
	// remove everything except digits
	re := regexp.MustCompile(`\D`)
	clean := re.ReplaceAllString(input, "")

	if clean == "" {
		return "", fmt.Errorf("invalid country code: %q", input)
	}

	return clean, nil
}
