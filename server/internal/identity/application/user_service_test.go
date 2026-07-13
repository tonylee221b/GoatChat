package application_test

import (
	"context"
	"errors"
	"testing"

	"GoatChat/GoatChat/internal/identity/application"
	"GoatChat/GoatChat/internal/identity/application/port/mocks"
	"GoatChat/GoatChat/internal/identity/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	for _, tt := range registerUserTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			err := svc.Register(context.Background(), tt.username)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFindUserByUsername(t *testing.T) {
	for _, tt := range findUserByUsernameTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			foundUser, err := svc.FindByUsername(context.Background(), tt.username)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, foundUser)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantUser, foundUser)
		})
	}
}

type registerUserTestCase struct {
	name      string
	username  domain.Username
	setupMock func(*mocks.MockUserRepository)
	wantErr   error
}

func registerUserTestCases() []registerUserTestCase {
	username, _ := domain.NewUsername("test-user")
	repoErr := errors.New(domain.ErrDB)

	return []registerUserTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
						return user.Username == username
					})).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "repository error",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(repoErr)
			},
			wantErr: repoErr,
		},
	}
}

type findUserByUsernameTestCase struct {
	name      string
	username  domain.Username
	setupMock func(*mocks.MockUserRepository)
	wantUser  *domain.User
	wantErr   error
}

func findUserByUsernameTestCases() []findUserByUsernameTestCase {
	username, _ := domain.NewUsername("test-user")
	user, _ := domain.NewUser(username)
	userNotFound := errors.New(domain.ErrUserNotFound)

	return []findUserByUsernameTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil)
			},
			wantUser: user,
			wantErr:  nil,
		},
		{
			name:     "user not found",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(nil, userNotFound)
			},
			wantUser: nil,
			wantErr:  userNotFound,
		},
	}
}
