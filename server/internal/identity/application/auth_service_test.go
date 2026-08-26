package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"GoatChat/GoatChat/internal/identity/application"
	"GoatChat/GoatChat/internal/identity/application/port/mocks"
	"GoatChat/GoatChat/internal/identity/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	for _, tt := range loginTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(
					mUserRepo,
					mSessionRepo,
					mPWHasher,
					mAccessTokenManager,
				)
			}

			svc := application.NewAuthService(
				testutils.StubTx{},
				mUserRepo,
				mSessionRepo,
				mPWHasher,
				mAccessTokenManager,
				testAccessTokenTTL,
				testRefreshTokenTTL,
			)

			tokenPair, err := svc.Login(context.Background(), tt.username, tt.plainPassword)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantAccessToken, tokenPair.AccessToken)
		})
	}
}

func TestLogout(t *testing.T) {
	for _, tt := range logoutTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mSessionRepo)
			}

			svc := application.NewAuthService(
				testutils.StubTx{},
				mUserRepo,
				mSessionRepo,
				mPWHasher,
				mAccessTokenManager,
				testAccessTokenTTL,
				testRefreshTokenTTL,
			)

			err := svc.Logout(context.Background(), tt.plainRefreshToken)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRefresh(t *testing.T) {
	for _, tt := range refreshTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mSessionRepo, mAccessTokenManager)
			}

			svc := application.NewAuthService(
				testutils.StubTx{},
				mUserRepo,
				mSessionRepo,
				mPWHasher,
				mAccessTokenManager,
				testAccessTokenTTL,
				testRefreshTokenTTL,
			)

			tokenPair, err := svc.Refresh(context.Background(), tt.plainRefreshToken)
			if tt.wantErrMsg != "" {
				require.ErrorContains(t, err, tt.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantAccessToken, tokenPair.AccessToken)
			require.NotEqual(t, tt.plainRefreshToken, tokenPair.RefreshToken)
			require.NotEqual(t, testRefreshTokenTTL, tokenPair.RefreshTokenExpiresAt)
		})
	}
}

const (
	testPlainPassword     = "plain-password"
	testPlainRefreshToken = "refresh-token"
	testAccessToken       = "access-token"
	testAccessTokenTTL    = 15 * time.Minute
	testRefreshTokenTTL   = 7 * 24 * time.Hour
)

type loginTestCase struct {
	name          string
	username      domain.Username
	plainPassword string
	setupMock     func(
		*mocks.MockUserRepository,
		*mocks.MockSessionRepository,
		*mocks.MockPasswordHasher,
		*mocks.MockAccessTokenManager,
	)
	wantAccessToken string
	wantErr         error
}

func loginTestCases() []loginTestCase {
	username, _ := domain.NewUsername("test-user")
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(username, passwordHash)
	invalidIDUser := domain.NewUser(username, passwordHash)
	invalidIDUser.ID = uuid.Nil

	userNotFound := errors.New(domain.ErrUserNotFound)
	passwordMismatch := errors.New(domain.ErrInvalidCredentials)
	saveSessionErr := errors.New(domain.ErrDB)
	sessionAlreadyExistsErr := errors.New("session already exists")
	invalidUserIDErr := errors.New(domain.ErrInvalidUserID)
	issueTokenErr := errors.New(domain.ErrInvalidConfig)

	return []loginTestCase{
		{
			name:          "success",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(nil).
					Once()

				sessionRepo.EXPECT().
					ExistsByUserID(mock.Anything, user.ID).
					Return(false, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(user.ID.String(), mock.Anything).
					Return(testAccessToken, nil).
					Once()

				sessionRepo.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(session domain.Session) bool {
						return session.UserID == user.ID &&
							!session.RefreshTokenHash.IsZero() &&
							session.ExpiresAt.After(session.CreatedAt)
					})).
					Return(nil).
					Once()
			},
			wantAccessToken: testAccessToken,
			wantErr:         nil,
		},
		{
			name:          "user not found",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockSessionRepository,
				_ *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(nil, userNotFound).
					Once()
			},
			wantErr: userNotFound,
		},
		{
			name:          "password mismatch",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(passwordMismatch).
					Once()
			},
			wantErr: passwordMismatch,
		},
		{
			name:          "invalid user ID",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(invalidIDUser, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(nil).
					Once()

				sessionRepo.EXPECT().
					ExistsByUserID(mock.Anything, invalidIDUser.ID).
					Return(false, nil).
					Once()
			},
			wantAccessToken: "",
			wantErr:         invalidUserIDErr,
		},
		{
			name:          "access token issue failure returns empty token pair",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(nil).
					Once()

				sessionRepo.EXPECT().
					ExistsByUserID(mock.Anything, user.ID).
					Return(false, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(user.ID.String(), mock.Anything).
					Return("", issueTokenErr).
					Once()
			},
			wantAccessToken: "",
			wantErr:         nil,
		},
		{
			name:          "session already exists",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(nil).
					Once()

				sessionRepo.EXPECT().
					ExistsByUserID(mock.Anything, user.ID).
					Return(true, nil).
					Once()
			},
			wantAccessToken: "",
			wantErr:         sessionAlreadyExistsErr,
		},
		{
			name:          "session save failure",
			username:      username,
			plainPassword: testPlainPassword,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, testPlainPassword).
					Return(nil).
					Once()

				sessionRepo.EXPECT().
					ExistsByUserID(mock.Anything, user.ID).
					Return(false, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(user.ID.String(), mock.Anything).
					Return(testAccessToken, nil).
					Once()

				sessionRepo.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(saveSessionErr).
					Once()
			},
			wantErr: saveSessionErr,
		},
	}
}

type logoutTestCase struct {
	name              string
	plainRefreshToken string
	setupMock         func(*mocks.MockSessionRepository)
	wantErr           error
}

func logoutTestCases() []logoutTestCase {
	refreshTokenHash, _ := domain.NewRefreshTokenHash(testPlainRefreshToken)
	now := time.Now().UTC()
	successSession, _ := domain.NewSession(
		uuid.New(),
		refreshTokenHash,
		now,
		now.Add(testRefreshTokenTTL),
	)
	revokeFailureSession, _ := domain.NewSession(
		uuid.New(),
		refreshTokenHash,
		now,
		now.Add(testRefreshTokenTTL),
	)

	sessionNotFound := errors.New(domain.ErrSessionNotFound)
	revokeSessionErr := errors.New(domain.ErrDB)

	return []logoutTestCase{
		{
			name:              "success",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(successSession, nil).
					Once()

				sessionRepo.EXPECT().
					Revoke(mock.Anything, successSession.ID, mock.Anything).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:              "empty refresh token",
			plainRefreshToken: "",
			setupMock:         nil,
			wantErr:           nil,
		},
		{
			name:              "session not found",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(nil, sessionNotFound).
					Once()
			},
			wantErr: sessionNotFound,
		},
		{
			name:              "session revoke failure",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(revokeFailureSession, nil).
					Once()

				sessionRepo.EXPECT().
					Revoke(mock.Anything, revokeFailureSession.ID, mock.Anything).
					Return(revokeSessionErr).
					Once()
			},
			wantErr: revokeSessionErr,
		},
	}
}

type refreshTestCase struct {
	name              string
	plainRefreshToken string
	setupMock         func(
		*mocks.MockSessionRepository,
		*mocks.MockAccessTokenManager,
	)
	wantAccessToken string
	wantErrMsg      string
}

func refreshTestCases() []refreshTestCase {
	currentHash, _ := domain.NewRefreshTokenHash(testPlainRefreshToken)
	now := time.Now().UTC()

	newSession := func(expiresAt time.Time, revokedAt *time.Time) *domain.Session {
		return domain.RestoreSession(
			uuid.New(),
			uuid.New(),
			currentHash,
			expiresAt,
			revokedAt,
			now.Add(-time.Hour),
		)
	}

	successSession := newSession(now.Add(testRefreshTokenTTL), nil)
	issueFailureSession := newSession(now.Add(testRefreshTokenTTL), nil)
	updateFailureSession := newSession(now.Add(testRefreshTokenTTL), nil)

	revokedAt := now.Add(-time.Minute)
	revokedSession := newSession(now.Add(testRefreshTokenTTL), &revokedAt)
	expiredSession := newSession(now.Add(-time.Minute), nil)

	return []refreshTestCase{
		{
			name:              "success",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(successSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(successSession.UserID.String(), successSession.ID.String()).
					Return(testAccessToken, nil).
					Once()

				sessionRepo.EXPECT().
					UpdateRefreshTokenHash(
						mock.Anything,
						successSession.ID,
						mock.MatchedBy(func(newHash domain.RefreshTokenHash) bool {
							return !newHash.IsZero() && !newHash.Equals(currentHash)
						}),
					).
					Return(nil).
					Once()
			},
			wantAccessToken: testAccessToken,
			wantErrMsg:      "",
		},
		{
			name:              "empty refresh token",
			plainRefreshToken: "",
			setupMock:         nil,
			wantAccessToken:   "",
			wantErrMsg:        domain.ErrEmptyRefreshToken,
		},
		{
			name:              "session not found",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(nil, errors.New(domain.ErrSessionNotFound)).
					Once()
			},
			wantAccessToken: "",
			wantErrMsg:      domain.ErrSessionNotFound,
		},
		{
			name:              "revoked session",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(revokedSession, nil).
					Once()
			},
			wantAccessToken: "",
			wantErrMsg:      domain.ErrSessionRevoked,
		},
		{
			name:              "expired session",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(expiredSession, nil).
					Once()
			},
			wantAccessToken: "",
			wantErrMsg:      domain.ErrSessionExpired,
		},
		{
			name:              "access token issue failure",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(issueFailureSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(issueFailureSession.UserID.String(), issueFailureSession.ID.String()).
					Return("", errors.New(domain.ErrTokenUserIDRequired)).
					Once()
			},
			wantAccessToken: "",
			wantErrMsg:      domain.ErrTokenUserIDRequired,
		},
		{
			name:              "refresh token update failure",
			plainRefreshToken: testPlainRefreshToken,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(updateFailureSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(updateFailureSession.UserID.String(), updateFailureSession.ID.String()).
					Return(testAccessToken, nil).
					Once()

				sessionRepo.EXPECT().
					UpdateRefreshTokenHash(
						mock.Anything,
						updateFailureSession.ID,
						mock.MatchedBy(func(newHash domain.RefreshTokenHash) bool {
							return !newHash.IsZero() && !newHash.Equals(currentHash)
						}),
					).
					Return(errors.New(domain.ErrSessionNotFound)).
					Once()
			},
			wantAccessToken: "",
			wantErrMsg:      domain.ErrSessionNotFound,
		},
	}
}
