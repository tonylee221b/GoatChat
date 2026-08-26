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
			mUserRepo := mocks.NewMockUserRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo, mPWHasher)
			}

			svc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			err := svc.Register(context.Background(), tt.username, tt.plainPassword)

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
			mUserRepo := mocks.NewMockUserRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo)
			}

			svc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
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
			mUserRepo := mocks.NewMockUserRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo)
			}

			svc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
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
	name          string
	username      domain.Username
	plainPassword string
	setupMock     func(
		*mocks.MockUserRepository,
		*mocks.MockPasswordHasher,
	)
	wantErr error
}

func registerUserTestCases() []registerUserTestCase {
	username, _ := domain.NewUsername("test-user")
	plainPassword := "plain-password"
	hashedPassword := "hashed-password"
	passwordHash, _ := domain.NewPasswordHash(hashedPassword)

	findUserErr := errors.New(domain.ErrDB)
	hashPasswordErr := errors.New("failed to hash password")
	saveUserErr := errors.New(domain.ErrDB)

	return []registerUserTestCase{
		{
			name:          "success",
			username:      username,
			plainPassword: plainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash(plainPassword).
					Return(hashedPassword, nil).
					Once()

				userRepo.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
						return user.Username == username &&
							user.PasswordHash == passwordHash
					})).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:          "user lookup failure",
			username:      username,
			plainPassword: plainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, findUserErr).
					Once()
			},
			wantErr: findUserErr,
		},
		{
			name:          "password hash failure",
			username:      username,
			plainPassword: plainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash(plainPassword).
					Return("", hashPasswordErr).
					Once()
			},
			wantErr: hashPasswordErr,
		},
		{
			name:          "user save failure",
			username:      username,
			plainPassword: plainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash(plainPassword).
					Return(hashedPassword, nil).
					Once()

				userRepo.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(saveUserErr).
					Once()
			},
			wantErr: saveUserErr,
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
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(username, passwordHash)
	userNotFound := errors.New(domain.ErrUserNotFound)

	return []findUserByUsernameTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
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
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
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
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	contact, _ := domain.NewContact("017-1234-2345", "gochat@example.com")

	successUser := domain.NewUser(username, passwordHash)
	wantSuccessUser := *successUser
	_ = wantSuccessUser.UpdateContact(contact)

	updateFailureUser := domain.NewUser(username, passwordHash)
	wantUpdateFailureUser := *updateFailureUser
	_ = wantUpdateFailureUser.UpdateContact(contact)

	userNotFound := errors.New(domain.ErrUserNotFound)
	updateUserErr := errors.New(domain.ErrDB)

	return []updateContactTestCase{
		{
			name:     "success",
			username: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(successUser, nil).
					Once()

				userRepo.EXPECT().
					Update(mock.Anything, wantSuccessUser).
					Return(nil).
					Once()
			},
			wantUser: &wantSuccessUser,
			wantErr:  nil,
		},
		{
			name:     "user not found",
			username: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(nil, userNotFound).
					Once()
			},
			wantUser: nil,
			wantErr:  userNotFound,
		},
		{
			name:     "user update failure",
			username: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(updateFailureUser, nil).
					Once()

				userRepo.EXPECT().
					Update(mock.Anything, wantUpdateFailureUser).
					Return(updateUserErr).
					Once()
			},
			wantUser: nil,
			wantErr:  updateUserErr,
		},
	}
}
