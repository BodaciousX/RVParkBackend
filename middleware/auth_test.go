// middleware/auth_test.go

package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BodaciousX/RVParkBackend/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock user service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user user.User, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func (m *MockUserService) GetUser(id string) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) Login(creds user.LoginCredentials) (*user.User, string, error) {
	args := m.Called(creds)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*user.User), args.String(1), args.Error(2)
}

func (m *MockUserService) ValidateToken(token string) (*user.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(userID string, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) RevokeAllTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestRequireAuth(t *testing.T) {
	// Setup
	mockUserService := new(MockUserService)
	authMiddleware := NewAuthMiddleware(mockUserService)

	// Create a simple handler for testing
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		user := r.Context().Value(userContextKey).(*user.User)
		assert.NotNil(t, user)
		assert.Equal(t, "user123", user.ID)
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := authMiddleware.RequireAuth(nextHandler)

	// Test case: no Authorization header
	req := httptest.NewRequest("GET", "http://example.com", nil)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Test case: invalid Authorization header format
	req = httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	recorder = httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Test case: invalid token
	req = httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	recorder = httptest.NewRecorder()

	mockUserService.On("ValidateToken", "invalidtoken").Return(nil, errors.New("invalid token")).Once()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	mockUserService.AssertExpectations(t)

	// Test case: valid token
	req = httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	recorder = httptest.NewRecorder()

	validUser := &user.User{
		ID:       "user123",
		Email:    "test@example.com",
		Username: "testuser",
		Role:     user.RoleStaff,
	}

	mockUserService.On("ValidateToken", "validtoken").Return(validUser, nil).Once()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestRequireAdmin(t *testing.T) {
	// Setup
	mockUserService := new(MockUserService)
	authMiddleware := NewAuthMiddleware(mockUserService)

	// Create a simple handler for testing
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := authMiddleware.RequireAdmin(nextHandler)

	// Test case: no user in context
	req := httptest.NewRequest("GET", "http://example.com", nil)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Test case: non-admin user
	req = httptest.NewRequest("GET", "http://example.com", nil)
	staffUser := &user.User{
		ID:       "user123",
		Email:    "staff@example.com",
		Username: "staffuser",
		Role:     user.RoleStaff,
	}

	ctx := context.WithValue(req.Context(), userContextKey, staffUser)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)

	// Test case: admin user
	req = httptest.NewRequest("GET", "http://example.com", nil)
	adminUser := &user.User{
		ID:       "admin123",
		Email:    "admin@example.com",
		Username: "adminuser",
		Role:     "ADMIN",
	}

	ctx = context.WithValue(req.Context(), userContextKey, adminUser)
	req = req.WithContext(ctx)
	recorder = httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRevokeUserTokens(t *testing.T) {
	// Setup
	mockUserService := new(MockUserService)
	authMiddleware := NewAuthMiddleware(mockUserService)

	// Test case: successful revocation
	mockUserService.On("RevokeAllTokens", "user123").Return(nil).Once()

	err := authMiddleware.RevokeUserTokens("user123")
	assert.NoError(t, err)
	mockUserService.AssertExpectations(t)

	// Test case: error
	mockUserService.On("RevokeAllTokens", "user123").Return(errors.New("db error")).Once()

	err = authMiddleware.RevokeUserTokens("user123")
	assert.Error(t, err)
	mockUserService.AssertExpectations(t)
}
