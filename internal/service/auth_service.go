package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/cache"
	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
	generateotp "github.com/alanzhumalin/bank/pkg/generate_otp"
	"github.com/alanzhumalin/bank/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	tokenKey       *string
	authRepository repository.AuthRepository
	txManager      repository.TxManagerRepository
	userService    UserService
	blackListToken cache.TokenBlackList
	otpRedis       cache.OTPStore
}

func NewAuthService(otpRedis cache.OTPStore, blackListToken cache.TokenBlackList, tokenKey *string, authRepository repository.AuthRepository, userService UserService, txManager repository.TxManagerRepository) AuthService {
	return &authService{otpRedis: otpRedis, blackListToken: blackListToken, tokenKey: tokenKey, authRepository: authRepository, userService: userService, txManager: txManager}
}

func (a *authService) Register(ctx context.Context, req dto.RegisterRequest, ip string, device string) (*dto.TokenPair, error) {

	token := &dto.TokenPair{}

	if err := a.txManager.WithTx(ctx, func(ctx context.Context) error {
		newUser := dto.CreateUserRequest{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Birthday:    *req.Birthday,
			PhoneNumber: req.PhoneNumber,
			Password:    req.Password,
		}

		id, role, err := a.userService.Create(ctx, newUser)
		fmt.Println(id, role)
		if err != nil {
			return err
		}

		sessionId := uuid.New().String()
		accessToken, err := jwt.GeneratateAccessToken(id, role, sessionId, *a.tokenKey)
		if err != nil {
			return fmt.Errorf("Error generating access token: %w", err)
		}
		refreshToken, expiresAt, createdAt, err := jwt.GeneratateRefreshToken(id, role, sessionId, *a.tokenKey)
		if err != nil {
			return fmt.Errorf("Error generating refresh token: %w", err)
		}

		hashedRefreshToken := jwt.HashRefreshToken(refreshToken, *a.tokenKey)

		session := &domain.Session{
			Id:                 sessionId,
			HashedRefreshToken: hashedRefreshToken,
			UserId:             id,
			Device:             device,
			Ip:                 ip,
			ExpiresAt:          *expiresAt,
			CreatedAt:          *createdAt,
		}

		err = a.authRepository.Сreate(ctx, *session)

		if err != nil {
			return err
		}

		token = &dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return token, nil
}

func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {

	userDetails, err := a.authRepository.GetDetails(ctx, req.PhoneNumber)

	if err != nil {
		return "", err
	}
	hashedPassword := userDetails.Password

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", domain.ErrorPasswordNotCorrect
	}

	challengeId := uuid.NewString()

	otp, err := generateotp.GenerateOtp()
	fmt.Println(otp)
	if err != nil {
		return "", err
	}
	hashedOtp, err := generateotp.HashOtp(otp)

	if err != nil {
		return "", err
	}

	otpDetail := domain.OTPDetail{
		UserId:      userDetails.UserId,
		Attempt:     1,
		PhoneNumber: req.PhoneNumber,
		CodeHash:    hashedOtp,
	}

	if err = a.otpRedis.Save(ctx, "login", challengeId, otpDetail); err != nil {
		return "", err
	}

	return challengeId, nil

}

func (a *authService) OTP(ctx context.Context, req dto.OTPRequest, ip string, device string) (*dto.TokenPair, error) {
	codeHash, err := generateotp.HashOtp(req.Code)
	if err != nil {
		return nil, err
	}
	ok, phoneNumber, err := a.otpRedis.Verify(ctx, "login", req.ChallengeId, codeHash)

	if err != nil {
		return nil, err
	}

	if !ok {
		return &dto.TokenPair{}, domain.ErrorOTPIncorrect
	}

	userDetails, err := a.authRepository.GetDetails(ctx, phoneNumber)

	sessionId := uuid.New().String()

	accessToken, err := jwt.GeneratateAccessToken(userDetails.UserId, userDetails.Role, sessionId, *a.tokenKey)
	if err != nil {
		return nil, fmt.Errorf("Error generating access token: %w", err)
	}
	refreshToken, expiresAt, createdAt, err := jwt.GeneratateRefreshToken(userDetails.UserId, userDetails.Role, sessionId, *a.tokenKey)
	if err != nil {
		return nil, fmt.Errorf("Error generating refresh token: %w", err)
	}

	hashedRefreshToken := jwt.HashRefreshToken(refreshToken, *a.tokenKey)

	session := &domain.Session{
		Id:                 sessionId,
		HashedRefreshToken: hashedRefreshToken,
		UserId:             userDetails.UserId,
		Device:             device,
		Ip:                 ip,
		ExpiresAt:          *expiresAt,
		CreatedAt:          *createdAt,
	}

	if err = a.authRepository.Сreate(ctx, *session); err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) UpdateSession(ctx context.Context, req dto.RefreshRequest) (*dto.TokenPair, string, error) {

	var token *dto.TokenPair
	var sId string

	if err := a.txManager.WithTx(ctx, func(ctx context.Context) error {

		claims, err := jwt.ParseAndValidateToken(req.RefreshToken, *a.tokenKey)

		if err != nil {
			return err
		}

		session, err := a.authRepository.GetSessionById(ctx, claims.SessionId)

		if err != nil {
			return err
		}

		if session.HashedRefreshToken != jwt.HashRefreshToken(req.RefreshToken, *a.tokenKey) {
			return domain.ErrorIncorrectRefreshToken
		}

		if session.ExpiresAt.Before(time.Now()) {
			return domain.ErrorRefreshTokenExpired
		}

		if !session.IsActive {
			return domain.ErrorSessionNotActive
		}

		accessToken, err := jwt.GeneratateAccessToken(claims.UserId, claims.Role, claims.SessionId, *a.tokenKey)

		if err != nil {
			return err
		}

		refreshToken, expiresAt, _, err := jwt.GeneratateRefreshToken(claims.UserId, claims.Role, claims.SessionId, *a.tokenKey)
		if err != nil {
			return err
		}

		hashedRefreshToken := jwt.HashRefreshToken(refreshToken, *a.tokenKey)

		if err = a.authRepository.Update(ctx, hashedRefreshToken, *expiresAt, claims.SessionId); err != nil {
			return err
		}

		token = &dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		sId = claims.SessionId
		return nil

	}); err != nil {
		return nil, "", err
	}

	return token, sId, nil

}

func (a *authService) Logout(ctx context.Context, sessionId string, jti string, exp time.Time) error {
	if err := a.authRepository.Revoke(ctx, sessionId); err != nil {
		return err
	}

	ttl := time.Until(exp)

	if ttl > 0 {
		if err := a.blackListToken.Add(ctx, jti, ttl+(10*time.Second)); err != nil {
			return err
		}
	}
	return nil
}

func (a *authService) LogoutFromAllDevices(ctx context.Context, userId int) error {
	return a.authRepository.RevokeAllUserDevices(ctx, userId)
}
