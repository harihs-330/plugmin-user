package usecase_test

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"user/adapter/apihook"
// 	"user/adapter/emailer"
// 	"user/config"
// 	"user/errortools"
// 	"user/internal/entities"
// 	"user/internal/repo/mock"
// 	"user/internal/usecase"

// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// // Mocked helper functions to set up test environment
// func settingupTest(t *testing.T) (*gomock.Controller, *mock.MockUserImply, usecase.UserImply, context.Context) {
// 	ctrl := gomock.NewController(t)
// 	mockRepo := mock.NewMockUserImply(ctrl)
// 	user := usecase.NewUser(mockRepo, &config.EnvConfig{}, &emailer.UnImplemented{}, apihook.HTTPAPI{})
// 	ctx := context.Background()

// 	return ctrl, mockRepo, user, ctx
// }

// // Success case: All permissions are listed successfully
// func TestListAllPermissions_Success(t *testing.T) {
// 	ctrl, mockRepo, user, ctx := settingupTest(t)
// 	defer ctrl.Finish()

// 	// Define expected permissions with PermissionDetail containing Name and Description
// 	expectedPermissions := &entities.Permissions{
// 		PermissionsMap: map[uuid.UUID]entities.PermissionDetail{
// 			uuid.MustParse("874dd293-84a6-4df6-8bab-5312b6dc4211"): {
// 				Name:        "read",
// 				Description: "Allows reading data",
// 			},
// 			uuid.MustParse("219947ba-fce9-4d9a-a497-02ec661f55d1"): {
// 				Name:        "write",
// 				Description: "Allows writing data",
// 			},
// 		},
// 	}

// 	// Mock the ListAllPermissions repository call to return the expected permissions
// 	mockRepo.EXPECT().
// 		ListAllPermissions(ctx).
// 		Return(expectedPermissions, nil).
// 		Times(1)

// 	// Call the method under test
// 	result, err := user.ListAllPermissions(ctx)
// 	if err != nil {
// 		require.NoError(t, err, fmt.Sprintf("expected no error, got %v", err))
// 	}
// 	// Validate the results
// 	require.NotNil(t, result, "expected non-nil result")
// 	require.Equal(t, len(expectedPermissions.PermissionsMap), len(result.PermissionsMap), "failed")
// 	require.Equal(t, expectedPermissions.PermissionsMap, result.PermissionsMap, "unexpected permission list")
// }

// // Failure case 1: No permissions exist
// func TestListAllPermissions_NoPermissions(t *testing.T) {
// 	ctrl, mockRepo, user, ctx := settingupTest(t)
// 	defer ctrl.Finish()

// 	mockRepo.EXPECT().
// 		ListAllPermissions(ctx).
// 		Return(nil, nil).
// 		Times(1)

// 	result, err := user.ListAllPermissions(ctx)
// 	if err != nil {
// 		require.NoError(t, err, fmt.Sprintf("expected no error, got %v", err))
// 	}
// 	require.Nil(t, result, "expected nil result when no permissions are found")
// }

// // Failure case 2: Repository returns an error
// func TestListAllPermissions_RepoError(t *testing.T) {
// 	ctrl, mockRepo, user, ctx := settingupTest(t)
// 	defer ctrl.Finish()

// 	mockRepo.EXPECT().
// 		ListAllPermissions(ctx).
// 		Return(nil, errors.New("database query failed")).
// 		Times(1)

// 	// Call the method under test
// 	result, err := user.ListAllPermissions(ctx)
// 	require.Nil(t, result, "expected nil result when repo fails")
// 	require.NotNil(t, err, "expected error when repo fails")
// 	require.NotEqual(t, errortools.FailedGettingPermission, err.Code, "unexpected error code")
// }
