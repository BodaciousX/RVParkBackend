// payment/p_service_test.go
package payment

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

func (m *MockRepository) Create(payment Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockRepository) Get(id string) (*Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payment), args.Error(1)
}

func (m *MockRepository) Update(payment Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) ListByTenant(tenantID string) ([]Payment, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]Payment), args.Error(1)
}

func (m *MockRepository) ListByDateRange(start, end time.Time) ([]Payment, error) {
	args := m.Called(start, end)
	return args.Get(0).([]Payment), args.Error(1)
}

func (m *MockRepository) ListByDateRangeAndTenant(start, end time.Time, tenantID string) ([]Payment, error) {
	args := m.Called(start, end, tenantID)
	return args.Get(0).([]Payment), args.Error(1)
}

func (m *MockRepository) GetLatestByTenant(tenantID string) (*Payment, error) {
	args := m.Called(tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payment), args.Error(1)
}

func TestCreatePayment_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	now := time.Now()
	nextPayment := now.AddDate(0, 1, 0)
	testPayment := Payment{
		TenantID:        uuid.New().String(),
		AmountDue:       500.00,
		DueDate:         now,
		NextPaymentDate: nextPayment,
	}

	// Setup expectations
	mockRepo.On("Create", mock.AnythingOfType("Payment")).Return(nil)

	// Call method being tested
	err := service.CreatePayment(testPayment)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify payment was created with the right values
	createCall := mockRepo.Calls[0]
	createdPayment := createCall.Arguments[0].(Payment)
	assert.Equal(t, testPayment.TenantID, createdPayment.TenantID)
	assert.Equal(t, testPayment.AmountDue, createdPayment.AmountDue)
	assert.Equal(t, testPayment.DueDate, createdPayment.DueDate)
	assert.Equal(t, testPayment.NextPaymentDate, createdPayment.NextPaymentDate)
	assert.NotEmpty(t, createdPayment.ID)
	assert.False(t, createdPayment.CreatedAt.IsZero())
	assert.False(t, createdPayment.UpdatedAt.IsZero())
}

func TestCreatePayment_ValidationFailure(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test cases for validation failures
	testCases := []struct {
		name    string
		payment Payment
		errMsg  string
	}{
		{
			name: "Empty tenant ID",
			payment: Payment{
				AmountDue:       500.00,
				DueDate:         time.Now(),
				NextPaymentDate: time.Now().AddDate(0, 1, 0),
			},
			errMsg: "tenant ID is required",
		},
		{
			name: "Zero amount due",
			payment: Payment{
				TenantID:        uuid.New().String(),
				AmountDue:       0,
				DueDate:         time.Now(),
				NextPaymentDate: time.Now().AddDate(0, 1, 0),
			},
			errMsg: "amount due must be greater than 0",
		},
		{
			name: "Empty due date",
			payment: Payment{
				TenantID:        uuid.New().String(),
				AmountDue:       500.00,
				NextPaymentDate: time.Now().AddDate(0, 1, 0),
			},
			errMsg: "due date is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call method being tested
			err := service.CreatePayment(tc.payment)

			// Assert expectations
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
			// Create should not be called for invalid payment
			mockRepo.AssertNotCalled(t, "Create", mock.Anything)
		})
	}
}

func TestGetPayment_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	now := time.Now()
	paymentID := uuid.New().String()
	testPayment := &Payment{
		ID:              paymentID,
		TenantID:        uuid.New().String(),
		AmountDue:       500.00,
		DueDate:         now,
		NextPaymentDate: now.AddDate(0, 1, 0),
		CreatedAt:       now.Add(-time.Hour),
		UpdatedAt:       now.Add(-time.Hour),
	}

	// Setup expectations
	mockRepo.On("Get", paymentID).Return(testPayment, nil)

	// Call method being tested
	payment, err := service.GetPayment(paymentID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, testPayment.ID, payment.ID)
	assert.Equal(t, testPayment.TenantID, payment.TenantID)
	assert.Equal(t, testPayment.AmountDue, payment.AmountDue)
	mockRepo.AssertExpectations(t)
}

func TestGetPayment_NotFound(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	paymentID := uuid.New().String()

	// Setup expectations
	mockRepo.On("Get", paymentID).Return(nil, errors.New("payment not found"))

	// Call method being tested
	payment, err := service.GetPayment(paymentID)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, payment)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePayment_Success(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	now := time.Now()
	paymentID := uuid.New().String()
	tenantID := uuid.New().String()
	paidDate := now.Add(-time.Hour)

	existingPayment := &Payment{
		ID:              paymentID,
		TenantID:        tenantID,
		AmountDue:       500.00,
		DueDate:         now,
		NextPaymentDate: now.AddDate(0, 1, 0),
		CreatedAt:       now.Add(-24 * time.Hour),
		UpdatedAt:       now.Add(-24 * time.Hour),
	}

	updatedPayment := Payment{
		ID:              paymentID,
		AmountDue:       450.00,               // Changed
		DueDate:         now.AddDate(0, 0, 7), // Changed
		PaidDate:        &paidDate,            // Added
		NextPaymentDate: now.AddDate(0, 1, 7), // Changed
	}

	// Setup expectations
	mockRepo.On("Get", paymentID).Return(existingPayment, nil)
	mockRepo.On("Update", mock.AnythingOfType("Payment")).Return(nil)

	// Call method being tested
	err := service.UpdatePayment(updatedPayment)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify payment was updated correctly
	updateCall := mockRepo.Calls[1]
	finalPayment := updateCall.Arguments[0].(Payment)
	assert.Equal(t, paymentID, finalPayment.ID)
	assert.Equal(t, tenantID, finalPayment.TenantID) // Should preserve original tenant ID
	assert.Equal(t, updatedPayment.AmountDue, finalPayment.AmountDue)
	assert.Equal(t, updatedPayment.DueDate, finalPayment.DueDate)
	assert.Equal(t, updatedPayment.PaidDate, finalPayment.PaidDate)
	assert.Equal(t, updatedPayment.NextPaymentDate, finalPayment.NextPaymentDate)
	assert.Equal(t, existingPayment.CreatedAt, finalPayment.CreatedAt)    // Should preserve created at
	assert.NotEqual(t, existingPayment.UpdatedAt, finalPayment.UpdatedAt) // Should update updated at
}

func TestUpdatePayment_ValidationFailure(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	now := time.Now()
	paymentID := uuid.New().String()
	existingPayment := &Payment{
		ID:              paymentID,
		TenantID:        uuid.New().String(),
		AmountDue:       500.00,
		DueDate:         now,
		NextPaymentDate: now.AddDate(0, 1, 0),
		CreatedAt:       now.Add(-24 * time.Hour),
		UpdatedAt:       now.Add(-24 * time.Hour),
	}

	// Test cases for validation failures
	testCases := []struct {
		name    string
		payment Payment
		errMsg  string
	}{
		{
			name: "Zero amount due",
			payment: Payment{
				ID:              paymentID,
				AmountDue:       0,
				DueDate:         now.AddDate(0, 0, 7),
				NextPaymentDate: now.AddDate(0, 1, 7),
			},
			errMsg: "amount due must be greater than 0",
		},
		{
			name: "Empty due date",
			payment: Payment{
				ID:              paymentID,
				AmountDue:       450.00,
				NextPaymentDate: now.AddDate(0, 1, 7),
			},
			errMsg: "due date is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup expectations for each test case
			mockRepo.On("Get", paymentID).Return(existingPayment, nil).Once()

			// Call method being tested
			err := service.UpdatePayment(tc.payment)

			// Assert expectations
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
			// Update should not be called for invalid payment
			mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		})
	}
}

func TestDeletePayment(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	paymentID := uuid.New().String()

	// Setup expectations
	mockRepo.On("Delete", paymentID).Return(nil)

	// Call method being tested
	err := service.DeletePayment(paymentID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTenantPayments(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	tenantID := uuid.New().String()
	testPayments := []Payment{
		{ID: uuid.New().String(), TenantID: tenantID, AmountDue: 500.00},
		{ID: uuid.New().String(), TenantID: tenantID, AmountDue: 600.00},
	}

	// Setup expectations - we should use ListByDateRangeAndTenant with a 6-month window
	mockRepo.On("ListByDateRangeAndTenant",
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"),
		tenantID).Return(testPayments, nil)

	// Call method being tested
	payments, err := service.GetTenantPayments(tenantID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, testPayments, payments)
	mockRepo.AssertExpectations(t)

	// Verify date range is approximately 6 months
	dateRangeCall := mockRepo.Calls[0]
	start := dateRangeCall.Arguments[0].(time.Time)
	end := dateRangeCall.Arguments[1].(time.Time)

	// Difference should be approximately 6 months (with a bit of tolerance for test execution time)
	diff := end.Sub(start)
	sixMonths := time.Hour * 24 * 30 * 6
	assert.InDelta(t, sixMonths.Hours(), diff.Hours(), 24) // Allow 1 day of difference
}

func TestGetPaymentsByDateRange(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	testPayments := []Payment{
		{ID: uuid.New().String(), TenantID: uuid.New().String(), AmountDue: 500.00, DueDate: start.AddDate(0, 0, 15)},
		{ID: uuid.New().String(), TenantID: uuid.New().String(), AmountDue: 600.00, DueDate: start.AddDate(0, 0, 25)},
	}

	// Setup expectations
	mockRepo.On("ListByDateRange", start, end).Return(testPayments, nil)

	// Call method being tested
	payments, err := service.GetPaymentsByDateRange(start, end)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, testPayments, payments)
	mockRepo.AssertExpectations(t)
}

func TestGetLatestPayment(t *testing.T) {
	// Create mock
	mockRepo := new(MockRepository)

	// Create service with mock
	service := NewService(mockRepo)

	// Test data
	tenantID := uuid.New().String()
	now := time.Now()

	testPayment := &Payment{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		AmountDue:       500.00,
		DueDate:         now,
		NextPaymentDate: now.AddDate(0, 1, 0),
	}

	// Setup expectations
	mockRepo.On("GetLatestByTenant", tenantID).Return(testPayment, nil)

	// Call method being tested
	payment, err := service.GetLatestPayment(tenantID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, testPayment, payment)
	mockRepo.AssertExpectations(t)
}
