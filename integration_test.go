package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BodaciousX/RVParkBackend/api"
	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services
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

type MockTenantService struct {
	mock.Mock
}

func (m *MockTenantService) ListTenants() ([]tenant.Tenant, error) {
	args := m.Called()
	return args.Get(0).([]tenant.Tenant), args.Error(1)
}

func (m *MockTenantService) CreateTenant(tenant tenant.Tenant) error {
	args := m.Called(tenant)
	return args.Error(0)
}

func (m *MockTenantService) GetTenant(id string) (*tenant.Tenant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.Tenant), args.Error(1)
}

func (m *MockTenantService) UpdateTenant(tenant tenant.Tenant) error {
	args := m.Called(tenant)
	return args.Error(0)
}

func (m *MockTenantService) DeleteTenant(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTenantService) GetTenantBySpace(spaceID string) (*tenant.Tenant, error) {
	args := m.Called(spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.Tenant), args.Error(1)
}

type MockSpaceService struct {
	mock.Mock
}

func (m *MockSpaceService) ListSpaces() (map[string][]space.Space, error) {
	args := m.Called()
	return args.Get(0).(map[string][]space.Space), args.Error(1)
}

func (m *MockSpaceService) GetSpace(id string) (*space.Space, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*space.Space), args.Error(1)
}

func (m *MockSpaceService) GetVacantSpaces() ([]space.Space, error) {
	args := m.Called()
	return args.Get(0).([]space.Space), args.Error(1)
}

func (m *MockSpaceService) ReserveSpace(spaceID string) error {
	args := m.Called(spaceID)
	return args.Error(0)
}

func (m *MockSpaceService) UnreserveSpace(spaceID string) error {
	args := m.Called(spaceID)
	return args.Error(0)
}

func (m *MockSpaceService) MoveIn(spaceID string, tenantID string) error {
	args := m.Called(spaceID, tenantID)
	return args.Error(0)
}

func (m *MockSpaceService) MoveOut(spaceID string) error {
	args := m.Called(spaceID)
	return args.Error(0)
}

func (m *MockSpaceService) UpdateSpace(space space.Space) error {
	args := m.Called(space)
	return args.Error(0)
}

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePayment(payment payment.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentService) GetPayment(id string) (*payment.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Payment), args.Error(1)
}

func (m *MockPaymentService) UpdatePayment(payment payment.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentService) DeletePayment(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentService) GetTenantPayments(tenantID string) ([]payment.Payment, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]payment.Payment), args.Error(1)
}

func (m *MockPaymentService) GetPaymentsByDateRange(start, end time.Time) ([]payment.Payment, error) {
	args := m.Called(start, end)
	return args.Get(0).([]payment.Payment), args.Error(1)
}

func (m *MockPaymentService) GetLatestPayment(tenantID string) (*payment.Payment, error) {
	args := m.Called(tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Payment), args.Error(1)
}

// Helper function to set up the server with mock services
func setupTestServer() (*api.Server, *MockUserService, *MockTenantService, *MockSpaceService, *MockPaymentService) {
	mockUserService := new(MockUserService)
	mockTenantService := new(MockTenantService)
	mockSpaceService := new(MockSpaceService)
	mockPaymentService := new(MockPaymentService)

	// Create auth middleware with mock user service
	authMiddleware := middleware.NewAuthMiddleware(mockUserService)

	// Create server with mock services
	server := api.NewServer(
		mockUserService,
		mockTenantService,
		mockSpaceService,
		mockPaymentService,
		authMiddleware,
	)

	return server, mockUserService, mockTenantService, mockSpaceService, mockPaymentService
}

// 1. Authentication Flow - Test login functionality
func TestLogin(t *testing.T) {
	// Setup test server with mock services
	server, mockUserService, _, _, _ := setupTestServer()

	// Mock data
	testUser := &user.User{
		ID:       uuid.New().String(),
		Email:    "test@example.com",
		Username: "testuser",
		Role:     user.RoleStaff,
	}
	testToken := "test-token-12345"

	// Setup login expectation
	mockUserService.On("Login", mock.MatchedBy(func(creds user.LoginCredentials) bool {
		return creds.Email == "test@example.com" && creds.Password == "password123"
	})).Return(testUser, testToken, nil)

	// Create login request
	loginReq := api.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.Mux.ServeHTTP(rr, req)

	// Check login response
	assert.Equal(t, http.StatusOK, rr.Code)

	var loginResp api.LoginResponse
	err := json.Unmarshal(rr.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, loginResp.User.ID)
	assert.Equal(t, testToken, loginResp.Token)

	// Assert expectations
	mockUserService.AssertExpectations(t)
}

// 2. Space Reservation - Test reserve space functionality
func TestReserveSpace(t *testing.T) {
	// Setup test server with mock services
	server, mockUserService, _, mockSpaceService, _ := setupTestServer()

	// Create test data
	spaceID := "A1"
	testUser := &user.User{
		ID:       uuid.New().String(),
		Email:    "staff@example.com",
		Username: "staff",
		Role:     user.RoleStaff,
	}

	// Setup expectations
	mockUserService.On("ValidateToken", "test-token").Return(testUser, nil)
	mockSpaceService.On("ReserveSpace", spaceID).Return(nil)

	// Create reserve request
	req, _ := http.NewRequest("POST", "/spaces/"+spaceID+"/reserve", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	server.Mux.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert expectations
	mockUserService.AssertExpectations(t)
	mockSpaceService.AssertExpectations(t)
}

// 3. Tenant Creation - Test create tenant functionality
func TestCreateTenant(t *testing.T) {
	// Setup test server with mock services
	server, mockUserService, mockTenantService, _, _ := setupTestServer()

	// Create test data
	testUser := &user.User{
		ID:       uuid.New().String(),
		Email:    "staff@example.com",
		Username: "staff",
		Role:     user.RoleStaff,
	}

	// Setup expectations
	mockUserService.On("ValidateToken", "test-token").Return(testUser, nil)
	mockTenantService.On("CreateTenant", mock.AnythingOfType("tenant.Tenant")).Return(nil)

	// Create tenant request
	moveInDate := time.Now()
	createTenantReq := api.CreateTenantRequest{
		Name:       "John Doe",
		MoveInDate: moveInDate,
		SpaceID:    "A1",
	}
	tenantBody, _ := json.Marshal(createTenantReq)
	req, _ := http.NewRequest("POST", "/tenants", bytes.NewBuffer(tenantBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	server.Mux.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	var respTenant tenant.Tenant
	err := json.Unmarshal(rr.Body.Bytes(), &respTenant)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", respTenant.Name)
	assert.Equal(t, "A1", respTenant.SpaceID)

	// Assert expectations
	mockUserService.AssertExpectations(t)
	mockTenantService.AssertExpectations(t)
}

// 4. Payment Creation - Test create payment functionality
func TestCreatePayment(t *testing.T) {
	// Setup test server with mock services
	server, mockUserService, _, _, mockPaymentService := setupTestServer()

	// Create test data
	tenantID := uuid.New().String()
	testUser := &user.User{
		ID:       uuid.New().String(),
		Email:    "staff@example.com",
		Username: "staff",
		Role:     user.RoleStaff,
	}

	// Setup expectations
	mockUserService.On("ValidateToken", "test-token").Return(testUser, nil)
	mockPaymentService.On("CreatePayment", mock.AnythingOfType("payment.Payment")).Return(nil)

	// Create payment request
	dueDate := time.Now().Add(24 * time.Hour)
	nextPaymentDate := time.Now().Add(30 * 24 * time.Hour)
	paidDate := time.Time{} // Zero time for unpaid payment

	paymentReq := api.CreatePaymentRequest{
		TenantID:        tenantID,
		AmountDue:       500.0,
		DueDate:         dueDate,
		NextPaymentDate: nextPaymentDate,
		PaidDate:        paidDate,
	}
	paymentBody, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(paymentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	server.Mux.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert expectations
	mockUserService.AssertExpectations(t)
	mockPaymentService.AssertExpectations(t)
}

// 5. List Vacant Spaces - Test get vacant spaces functionality
func TestGetVacantSpaces(t *testing.T) {
	// Setup test server with mock services
	server, mockUserService, _, mockSpaceService, _ := setupTestServer()

	// Create test data
	testUser := &user.User{
		ID:       uuid.New().String(),
		Email:    "staff@example.com",
		Username: "staff",
		Role:     user.RoleStaff,
	}

	// Setup vacant spaces to return
	vacantSpaces := []space.Space{
		{
			ID:       "A1",
			Section:  "Mane Street",
			Status:   space.StatusVacant,
			Reserved: false,
		},
		{
			ID:       "B2",
			Section:  "Grace Street",
			Status:   space.StatusVacant,
			Reserved: false,
		},
	}

	// Setup expectations
	mockUserService.On("ValidateToken", "test-token").Return(testUser, nil)
	mockSpaceService.On("GetVacantSpaces").Return(vacantSpaces, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/spaces/vacant", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	rr := httptest.NewRecorder()
	server.Mux.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var respSpaces []space.Space
	err := json.Unmarshal(rr.Body.Bytes(), &respSpaces)
	assert.NoError(t, err)
	assert.Len(t, respSpaces, 2)
	assert.Equal(t, "A1", respSpaces[0].ID)
	assert.Equal(t, "B2", respSpaces[1].ID)

	// Assert expectations
	mockUserService.AssertExpectations(t)
	mockSpaceService.AssertExpectations(t)
}
