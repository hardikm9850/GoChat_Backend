package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	auth "github.com/hardikm9850/GoChat/internal/auth"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	internalDomain "github.com/hardikm9850/GoChat/internal/domain"
	authkitjwt "github.com/hardikm9850/authkit/jwt"
	authkitpassword "github.com/hardikm9850/authkit/password"
	"log"
	"time"
)

type authService struct {
	userRepo   repository.UserRepository
	refresh    repository.RefreshTokenRepository
	jwtManager authkitjwt.HS256Manager
}

func (s *authService) Logout(userID string) error {
	return s.refresh.DeleteByUserID(userID)
}

func New(
	userRepo repository.UserRepository,
	refresh repository.RefreshTokenRepository,
	jwtManager authkitjwt.HS256Manager,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		refresh:    refresh,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(countryCode, phone, password, name string) (Tokens, error) {
	log.Println("Service /auth/register hit")

	countryCode, err := internalDomain.NormalizeCountryCode(countryCode)
	if err != nil {
		return Tokens{}, err
	}

	// check if user exists
	_, err = s.userRepo.FindByMobile(phone, countryCode)
	log.Println("FindByMobile initial")
	if err == nil {
		return Tokens{}, auth.ErrUserAlreadyExists
	}

	// hash pwd
	hashedPassword, err := authkitpassword.HashPassword(password)
	if err != nil {
		return Tokens{}, err
	}

	phoneHash := internalDomain.HashPhoneNumber(countryCode, phone)

	user := domain.User{
		ID:           uuid.NewString(),
		Name:         name,
		PhoneNumber:  phone,
		PasswordHash: hashedPassword,
		PhoneHash:    phoneHash,
		CreatedAt:    time.Now(),
		CountryCode:  countryCode,
	}

	err = s.userRepo.Create(user)
	log.Println("Error creating user")
	if err != nil {
		return Tokens{}, err
	}
	log.Println("FindByMobile after creating user")
	user, err = s.userRepo.FindByMobile(phone, countryCode)
	if err != nil {
		log.Printf("FindByMobile error is %s\n", err)
		return Tokens{}, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate jwt: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

func (s *authService) Login(phone, password, countryCode string) (Tokens, error) {
	user, err := s.userRepo.FindByMobile(phone, countryCode)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return Tokens{}, auth.ErrUserDoesNotExists
		}
		return Tokens{}, err
	}

	if err := authkitpassword.VerifyPassword(password, user.PasswordHash); err != nil {
		return Tokens{}, auth.ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate jwt: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}

	refresh := domain.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt: time.Now(),
	}
	// Delete previous token
	_ = s.refresh.DeleteByUserID(user.ID)

	// Save refresh token to repo
	if err := s.refresh.Create(refresh); err != nil {
		return Tokens{}, fmt.Errorf("save refresh token: %w", err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

func (s *authService) RefreshAccessToken(refreshToken string) (Tokens, error) {
	// 1. Verify JWT itself
	claims, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, auth.ErrInvalidToken
	}

	// 2. Fetch storedRefreshToken refresh token
	storedRefreshToken, err := s.refresh.FindByUserID(claims.UserID)
	if err != nil {
		return Tokens{}, auth.ErrInvalidToken
	}

	// 3. Detect reuse
	if storedRefreshToken.Token != refreshToken {
		// possible token theft
		_ = s.refresh.DeleteByUserID(claims.UserID)
		return Tokens{}, auth.ErrTokenReuseDetected
	}

	// 4. Expiry check
	if time.Now().After(storedRefreshToken.ExpiresAt) {
		_ = s.refresh.DeleteByUserID(claims.UserID)
		return Tokens{}, auth.ErrExpiredToken
	}

	// 5. Generate new tokens
	newAccessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID)
	if err != nil {
		return Tokens{}, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return Tokens{}, err
	}

	// 6. Rotate refresh token (replace)
	err = s.refresh.UpdateByUserID(
		claims.UserID,
		newRefreshToken,
		time.Now().Add(30*24*time.Hour),
	)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
