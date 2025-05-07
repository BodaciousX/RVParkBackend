// user/u_service_test.go
package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) Get(id string) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetByEmail(email string) (*User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockTokenRepository is a mock implementation of the TokenRepository interface
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) CreateToken(token Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetToken(tokenHash string) (*Token, error) {
	args := m.Called(tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Token), args.Error(1)
}

func (m *MockTokenRepository) RevokeToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockTokenRepository) RevokeAllUserTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockTokenRepository) CleanExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTokenRepo := new(MockTokenRepository)

	// Create service with mocks
	service := NewService(mockRepo, mockTokenRepo)

	// Test data
	testUser := User{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     RoleStaff,
	}
	testPassword := "password123"

	// Setup expectations
	mockRepo.On("Create", mock.AnythingOfType("User")).Return(nil)

	// Call method being tested
	err := service.CreateUser(testUser, testPassword)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that password was hashed (indirectly)
	createCall := mockRepo.Calls[0]
	createdUser := createCall.Arguments[0].(User)
	assert.NotEmpty(t, createdUser.PasswordHash)
	assert.NotEqual(t, testPassword, createdUser.PasswordHash)
}

func TestLogin_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTokenRepo := new(MockTokenRepository)

	// Create service with mocks
	service := NewService(mockRepo, mockTokenRepo)

	// Hash a known password for our test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// Test data
	testUser := &User{
		ID:           "user123",
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         RoleStaff,
	}

	// Setup expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(testUser, nil)
	mockTokenRepo.On("CreateToken", mock.AnythingOfType("Token")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("User")).Return(nil)

	// Call method being tested
	user, token, err := service.Login(LoginCredentials{
		Email:    "test@example.com",
		Password: "correctpassword",
	})

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, testUser.ID, user.ID)
	mockRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTokenRepo := new(MockTokenRepository)

	// Create service with mocks
	service := NewService(mockRepo, mockTokenRepo)

	// Hash a known password for our test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// Test data
	testUser := &User{
		ID:           "user123",
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Role:         RoleStaff,
	}

	// Setup expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(testUser, nil)

	// Call method being tested with wrong password
	user, token, err := service.Login(LoginCredentials{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
	// Token repo methods should not be called
	mockTokenRepo.AssertNotCalled(t, "CreateToken", mock.Anything)
}

func TestValidateToken_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTokenRepo := new(MockTokenRepository)

	// Create service with mocks
	service := NewService(mockRepo, mockTokenRepo)

	// Test data
	testToken := "validtoken123"
	testTokenHash := "hashed_validtoken123" // Simplified for testing
	testUser := &User{ID: "user123", Email: "test@example.com"}

	// Create a token that is not expired and not revoked
	storedToken := &Token{
		TokenHash: testTokenHash,
		UserID:    testUser.ID,
		ExpiresAt: time.Now().Add(time.Hour), // Not expired
		Revoked:   false,
	}

	// Setup expectations
	mockTokenRepo.On("GetToken", mock.AnythingOfType("string")).Return(storedToken, nil)
	mockRepo.On("Get", testUser.ID).Return(testUser, nil)

	// Call method being tested
	user, err := service.ValidateToken(testToken)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.ID, user.ID)
	mockTokenRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestValidateToken_Expired(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTokenRepo := new(MockTokenRepository)

	// Create service with mocks
	service := NewService(mockRepo, mockTokenRepo)

	// Test data
	testToken := "expiredtoken123"
	testTokenHash := "hashed_expiredtoken123" // Simplified for testing

	// Create an expired token
	storedToken := &Token{
		TokenHash: testTokenHash,
		UserID:    "user123",
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
		Revoked:   false,
	}

	// Setup expectations
	mockTokenRepo.On("GetToken", mock.AnythingOfType("string")).Return(storedToken, nil)

	// Call method being tested
	user, err := service.ValidateToken(testToken)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	mockTokenRepo.AssertExpectations(t)
	// User repo should not be called for expired token
	mockRepo.AssertNotCalled(t, "Get", mock.Anything)
}
