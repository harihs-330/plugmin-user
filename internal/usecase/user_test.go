package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"user/adapter/apihook"
	"user/adapter/emailer"
	"user/config"
	"user/internal/entities"
	mockrepo "user/internal/repo/mock"
	"user/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) ListUser(ctx context.Context, pagination entities.Pagination,
	filter entities.UserFilter) ([]*entities.User, int64, error) {

	args := m.Called(ctx, pagination, filter)
	return args.Get(0).([]*entities.User), args.Get(1).(int64), args.Error(2) //nolint
}

func TestListUser(t *testing.T) { //nolint

	mockRepo := new(MockUserRepo)
	tests := []struct {
		name       string
		pagination entities.Pagination
		filter     entities.UserFilter
		mockSetup  func()
		wantErr    bool
	}{
		{
			name: "successful query",
			pagination: entities.Pagination{
				Page:  1,
				Limit: 10,
				Sort:  "created_on",
				Order: "ASC",
			},
			filter: entities.UserFilter{
				ID:           "",
				Name:         "",
				Email:        "",
				Organization: "",
				ProjectID:    uuid.New().String(),
			},
			mockSetup: func() {
				users := []*entities.User{
					{
						UserID:       "a1885813-24f9-4fee-aca4-0cb738510a8f",
						Name:         "John Doe",
						Email:        "john.doe@example.com",
						Purpose:      "student",
						ISActive:     true,
						Organization: "Org1",
						CreatedOn:    time.Now(),
						UpdatedOn:    time.Now(),
					},
				}
				mockRepo.On("ListUser", mock.Anything, mock.Anything, mock.Anything).Return(users, int64(1), nil)
			},
			wantErr: false,
		},
		{
			name: "query error",
			pagination: entities.Pagination{
				Page:  1,
				Limit: 10,
				Sort:  "created_on",
				Order: "ASC",
			},
			filter: entities.UserFilter{
				ID:           "",
				Name:         "",
				Email:        "",
				Organization: "",
				ProjectID:    uuid.New().String(),
			},
			mockSetup: func() {
				mockRepo.On("ListUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, int64(0), sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			users, total, err := mockRepo.ListUser(context.Background(), tt.pagination, tt.filter)
			if err != nil {
				t.Errorf("ListUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (len(users) == 0 || total == 0) {
				t.Errorf("ListUser() got no users or total count = %v", total)
			}
		})
	}

	mockRepo.AssertExpectations(t)
}

// Mocked helper functions to set up test environment
func settingupTest(t *testing.T) (*gomock.Controller, *mockrepo.MockUserImply, usecase.UserImply, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mockrepo.NewMockUserImply(ctrl)
	user := usecase.NewUser(mockRepo, &config.EnvConfig{}, &emailer.UnImplemented{}, apihook.HTTPAPI{})
	ctx := context.Background()

	return ctrl, mockRepo, user, ctx
}

// Success case: ListUserPermissions returns permissions successfully
func TestListUserPermissions_Success(t *testing.T) {
	ctrl, mockRepo, user, ctx := settingupTest(t)
	defer ctrl.Finish()

	userID := "e408c0f0-1a8c-4072-b2d3-e2e1d952b667"
	projectID := "b1761268-c77f-4192-944a-70d914de582f"

	expectedPermissions := []entities.PermissionDetail{
		{ID: "permission-1", Name: "read"},
		{ID: "permission-2", Name: "write"},
	}

	mockRepo.EXPECT().
		ListUserPermissions(ctx, userID, projectID).
		Return(expectedPermissions, nil).
		Times(1)

	result, err := user.ListUserPermissions(ctx, userID, projectID)

	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	require.NotNil(t, result, "expected non-nil result but got nil")
	require.Equal(t, expectedPermissions, result, "unexpected permissions list")
}

// Failure case: No permissions exist
func TestListUserPermissions_NoPermissions(t *testing.T) {
	ctrl, mockRepo, user, ctx := settingupTest(t)
	defer ctrl.Finish()

	userID := "e408c0f0-1a8c-4072-b2d3-e2e1d952b667"
	projectID := "b1761268-c77f-4192-944a-70d914de582f"

	mockRepo.EXPECT().
		ListUserPermissions(ctx, userID, projectID).
		Return(nil, nil).
		Times(1)

	result, err := user.ListUserPermissions(ctx, userID, projectID)

	require.Nil(t, err)
	require.Len(t, result, 0)
}
