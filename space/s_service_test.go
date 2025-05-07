// space/s_service_test.go
package space

import (
	"testing"

	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List() ([]Space, error) {
	args := m.Called()
	return args.Get(0).([]Space), args.Error(1)
}

func (m *MockRepository) Get(id string) (*Space, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Space), args.Error(1)
}

func (m *MockRepository) Update(space Space) error {
	args := m.Called(space)
	return args.Error(0)
}

// MockTenantService is a mock implementation of the tenant.Service interface
type MockTenantService struct {
	mock.Mock
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

func (m *MockTenantService) GetTenantBySpace(spaceID string) (*tenant.Tenant, error) {
	args := m.Called(spaceID)
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

func (m *MockTenantService) ListTenants() ([]tenant.Tenant, error) {
	args := m.Called()
	return args.Get(0).([]tenant.Tenant), args.Error(1)
}

func TestListSpaces(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data
	testSpaces := []Space{
		{ID: "A1", Section: "Mane Street", Status: StatusVacant},
		{ID: "A2", Section: "Mane Street", Status: StatusOccupied},
		{ID: "B1", Section: "Grace Street", Status: StatusVacant},
	}

	// Setup expectations
	mockRepo.On("List").Return(testSpaces, nil)

	// Call method being tested
	result, err := service.ListSpaces()

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2) // Two sections
	assert.Len(t, result["Mane Street"], 2)
	assert.Len(t, result["Grace Street"], 1)
	mockRepo.AssertExpectations(t)
}

func TestGetVacantSpaces(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data with some occupied, some vacant
	tenantID := "tenant1"
	testSpaces := []Space{
		{ID: "A1", Section: "Mane Street", Status: StatusVacant, Reserved: false},
		{ID: "A2", Section: "Mane Street", Status: StatusOccupied, TenantID: &tenantID},
		{ID: "B1", Section: "Grace Street", Status: StatusVacant, Reserved: false},
		{ID: "B2", Section: "Grace Street", Status: StatusReserved, Reserved: true},
	}

	// Setup expectations
	mockRepo.On("List").Return(testSpaces, nil)

	// Call method being tested
	vacantSpaces, err := service.GetVacantSpaces()

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, vacantSpaces, 2) // Only 2 spaces are vacant and not reserved
	assert.Equal(t, "A1", vacantSpaces[0].ID)
	assert.Equal(t, "B1", vacantSpaces[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestReserveSpace_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data
	spaceID := "A1"
	testSpace := &Space{
		ID:       spaceID,
		Section:  "Mane Street",
		Status:   StatusVacant,
		Reserved: false,
	}

	// Setup expectations
	mockRepo.On("Get", spaceID).Return(testSpace, nil)
	mockRepo.On("Update", mock.AnythingOfType("Space")).Return(nil)

	// Call method being tested
	err := service.ReserveSpace(spaceID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify space was updated correctly
	updateCall := mockRepo.Calls[1]
	updatedSpace := updateCall.Arguments[0].(Space)
	assert.Equal(t, spaceID, updatedSpace.ID)
	assert.Equal(t, StatusReserved, updatedSpace.Status)
	assert.True(t, updatedSpace.Reserved)
}

func TestReserveSpace_AlreadyOccupied(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data - space is already occupied
	spaceID := "A1"
	tenantID := "tenant1"
	testSpace := &Space{
		ID:       spaceID,
		Section:  "Mane Street",
		Status:   StatusOccupied,
		TenantID: &tenantID,
	}

	// Setup expectations
	mockRepo.On("Get", spaceID).Return(testSpace, nil)

	// Call method being tested
	err := service.ReserveSpace(spaceID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not vacant")
	mockRepo.AssertExpectations(t)
	// Update should not be called for occupied space
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestMoveIn_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data
	spaceID := "A1"
	tenantID := "tenant1"
	testSpace := &Space{
		ID:       spaceID,
		Section:  "Mane Street",
		Status:   StatusVacant,
		Reserved: false,
	}

	// Setup expectations
	mockRepo.On("Get", spaceID).Return(testSpace, nil)
	mockRepo.On("Update", mock.AnythingOfType("Space")).Return(nil)

	// Call method being tested
	err := service.MoveIn(spaceID, tenantID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify space was updated correctly
	updateCall := mockRepo.Calls[1]
	updatedSpace := updateCall.Arguments[0].(Space)
	assert.Equal(t, spaceID, updatedSpace.ID)
	assert.Equal(t, StatusOccupied, updatedSpace.Status)
	assert.Equal(t, tenantID, *updatedSpace.TenantID)
	assert.False(t, updatedSpace.Reserved)
}

func TestMoveOut_Success(t *testing.T) {
	// Create mocks
	mockRepo := new(MockRepository)
	mockTenantService := new(MockTenantService)

	// Create service with mocks
	service := NewService(mockRepo, mockTenantService)

	// Test data
	spaceID := "A1"
	tenantID := "tenant1"
	testSpace := &Space{
		ID:       spaceID,
		Section:  "Mane Street",
		Status:   StatusOccupied,
		TenantID: &tenantID,
	}

	// Setup expectations
	mockRepo.On("Get", spaceID).Return(testSpace, nil)
	mockRepo.On("Update", mock.AnythingOfType("Space")).Return(nil)

	// Call method being tested
	err := service.MoveOut(spaceID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify space was updated correctly
	updateCall := mockRepo.Calls[1]
	updatedSpace := updateCall.Arguments[0].(Space)
	assert.Equal(t, spaceID, updatedSpace.ID)
	assert.Equal(t, StatusVacant, updatedSpace.Status)
	assert.Nil(t, updatedSpace.TenantID)
	assert.False(t, updatedSpace.Reserved)
}
