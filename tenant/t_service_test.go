// tenant/t_service_test.go
package tenant

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List() ([]Tenant, error) {
	args := m.Called()
	return args.Get(0).([]Tenant), args.Error(1)
}

func (m *MockRepository) Create(tenant Tenant) error {
	args := m.Called(tenant)
	return args.Error(0)
}

func (m *MockRepository) Get(id string) (*Tenant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Tenant), args.Error(1)
}

func (m *MockRepository) GetBySpace(spaceID string) (*Tenant, error) {
	args := m.Called(spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Tenant), args.Error(1)
}

func (m *MockRepository) Update(tenant Tenant) error {
	args := m.Called(tenant)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTenant_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	testTenant := Tenant{
		ID:      uuid.New().String(),
		Name:    "John Doe",
		SpaceID: "A1",
	}

	// Setup expectations
	mockRepo.On("GetBySpace", testTenant.SpaceID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("Tenant")).Return(nil)

	// Call method being tested
	err := service.CreateTenant(testTenant)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify tenant was created with the right values
	createCall := mockRepo.Calls[1]
	createdTenant := createCall.Arguments[0].(Tenant)
	assert.Equal(t, testTenant.ID, createdTenant.ID)
	assert.Equal(t, testTenant.Name, createdTenant.Name)
	assert.Equal(t, testTenant.SpaceID, createdTenant.SpaceID)
	assert.False(t, createdTenant.MoveInDate.IsZero())
}

func TestCreateTenant_SpaceOccupied(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	testTenant := Tenant{
		ID:      uuid.New().String(),
		Name:    "John Doe",
		SpaceID: "A1",
	}
	existingTenant := &Tenant{
		ID:      uuid.New().String(),
		Name:    "Existing Tenant",
		SpaceID: "A1",
	}

	// Setup expectations - space is already occupied
	mockRepo.On("GetBySpace", testTenant.SpaceID).Return(existingTenant, nil)

	// Call method being tested
	err := service.CreateTenant(testTenant)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already occupied")
	mockRepo.AssertExpectations(t)
	// Create should not be called if space is occupied
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestUpdateTenant_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	now := time.Now()
	tenantID := uuid.New().String()
	oldSpaceID := "A1"
	newSpaceID := "B1"

	existingTenant := &Tenant{
		ID:         tenantID,
		Name:       "John Doe",
		SpaceID:    oldSpaceID,
		MoveInDate: now.Add(-24 * time.Hour),
		CreatedAt:  now.Add(-24 * time.Hour),
	}

	updatedTenant := Tenant{
		ID:      tenantID,
		Name:    "John Doe Updated",
		SpaceID: newSpaceID,
	}

	// Setup expectations
	mockRepo.On("Get", tenantID).Return(existingTenant, nil)
	mockRepo.On("GetBySpace", newSpaceID).Return(nil, errors.New("not found"))
	mockRepo.On("Update", mock.AnythingOfType("Tenant")).Return(nil)

	// Call method being tested
	err := service.UpdateTenant(updatedTenant)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify tenant was updated correctly
	updateCall := mockRepo.Calls[2]
	finalTenant := updateCall.Arguments[0].(Tenant)
	assert.Equal(t, updatedTenant.ID, finalTenant.ID)
	assert.Equal(t, updatedTenant.Name, finalTenant.Name)
	assert.Equal(t, updatedTenant.SpaceID, finalTenant.SpaceID)
	assert.Equal(t, existingTenant.MoveInDate, finalTenant.MoveInDate)
	assert.Equal(t, existingTenant.CreatedAt, finalTenant.CreatedAt)
}

func TestUpdateTenant_NewSpaceOccupied(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	tenantID := uuid.New().String()
	oldSpaceID := "A1"
	newSpaceID := "B1"

	existingTenant := &Tenant{
		ID:      tenantID,
		Name:    "John Doe",
		SpaceID: oldSpaceID,
	}

	updatedTenant := Tenant{
		ID:      tenantID,
		Name:    "John Doe Updated",
		SpaceID: newSpaceID,
	}

	existingTenantInNewSpace := &Tenant{
		ID:      uuid.New().String(),
		Name:    "Other Tenant",
		SpaceID: newSpaceID,
	}

	// Setup expectations
	mockRepo.On("Get", tenantID).Return(existingTenant, nil)
	mockRepo.On("GetBySpace", newSpaceID).Return(existingTenantInNewSpace, nil)

	// Call method being tested
	err := service.UpdateTenant(updatedTenant)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already occupied")
	mockRepo.AssertExpectations(t)
	// Update should not be called if new space is occupied
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestDeleteTenant_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	tenantID := uuid.New().String()
	existingTenant := &Tenant{
		ID:      tenantID,
		Name:    "John Doe",
		SpaceID: "A1",
	}

	// Setup expectations
	mockRepo.On("Get", tenantID).Return(existingTenant, nil)
	mockRepo.On("Delete", tenantID).Return(nil)

	// Call method being tested
	err := service.DeleteTenant(tenantID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTenant_NotFound(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	tenantID := uuid.New().String()

	// Setup expectations
	mockRepo.On("Get", tenantID).Return(nil, errors.New("tenant not found"))

	// Call method being tested
	err := service.DeleteTenant(tenantID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant not found")
	mockRepo.AssertExpectations(t)
	// Delete should not be called if tenant doesn't exist
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
}
