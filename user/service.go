// user/service.go contains the implementation of the user service.
package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo        Repository
	tokenRepo   TokenRepository
	tokenExpiry time.Duration
}

func NewService(repo Repository, tokenRepo TokenRepository) Service {
	return &service{
		repo:        repo,
		tokenRepo:   tokenRepo,
		tokenExpiry: 24 * time.Hour,
	}
}

// Added missing methods:

func (s *service) CreateUser(user User, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set user fields
	user.PasswordHash = string(hashedPassword)
	user.CreatedAt = time.Now()

	return s.repo.Create(user)
}

func (s *service) GetUser(id string) (*User, error) {
	return s.repo.Get(id)
}

func (s *service) UpdateUser(user User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id string) error {
	// Revoke all tokens for the user before deletion
	if err := s.tokenRepo.RevokeAllUserTokens(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *service) ChangePassword(userID string, oldPassword, newPassword string) error {
	// Get the user
	user, err := s.repo.Get(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(oldPassword),
	); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user with new password
	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(*user)
}

// Existing methods:

func (s *service) Login(creds LoginCredentials) (*User, string, error) {
	user, err := s.repo.GetByEmail(creds.Email)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(creds.Password),
	); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate new token
	token, tokenHash, err := GenerateToken()
	if err != nil {
		return nil, "", err
	}

	// Store token
	if err := s.tokenRepo.CreateToken(Token{
		TokenHash: tokenHash,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, "", err
	}

	// Update last login
	user.LastLogin = time.Now()
	if err := s.repo.Update(*user); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *service) ValidateToken(token string) (*User, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.tokenRepo.GetToken(tokenHash)
	if err != nil {
		return nil, err
	}

	if storedToken.ExpiresAt.Before(time.Now()) || storedToken.Revoked {
		return nil, errors.New("token is expired or revoked")
	}

	return s.repo.Get(storedToken.UserID)
}

// Helper function to generate tokens
func GenerateToken() (string, string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}
