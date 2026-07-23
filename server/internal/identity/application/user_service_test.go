package application_test

import (
	"context"
	"errors"
	"log/slog"
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
				return
			}

			require.NoError(t, err)
			require.Equal(t, *tt.wantUser, foundUser)
		})
	}
}

func TestUpdateContact(t *testing.T) {
	contact, _ := domain.NewContact("017-1234-2345", "gochat@example.com")

	for _, tt := range updateContactTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			updatedUser, err := svc.UpdateContact(context.Background(), tt.username, contact)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			slog.Info("want_user", "user", tt.wantUser)

			require.NoError(t, err)
			require.Equal(t, *tt.wantUser, updatedUser)
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
					ExistsByUsername(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				m.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
						return user.Username == username
					})).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:     "repository error",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					ExistsByUsername(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				m.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(repoErr).
					Once()
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
	user := domain.NewUser(username)
	userNotFound := errors.New(domain.ErrUserNotFound)

	return []findUserByUsernameTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()
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
					Return(nil, userNotFound).
					Once()
			},
			wantUser: nil,
			wantErr:  userNotFound,
		},
	}
}

type updateContactTestCase struct {
	name      string
	username  domain.Username
	setupMock func(*mocks.MockUserRepository)
	wantUser  *domain.User
	wantErr   error
}

func updateContactTestCases() []updateContactTestCase {
	username, _ := domain.NewUsername("test-user")
	user := domain.NewUser(username)
	userNotFound := errors.New(domain.ErrUserNotFound)

	return []updateContactTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(m *mocks.MockUserRepository) {
				user.Contact.PhoneNumber.Value = "017-1234-2345"
				user.Contact.Email.Value = "gochat@example.com"

				m.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(user, nil).
					Once()

				m.EXPECT().
					Update(mock.Anything, *user).
					Return(nil).
					Once()
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
					Return(nil, userNotFound).
					Once()
			},
			wantUser: nil,
			wantErr:  userNotFound,
		},
	}
}
