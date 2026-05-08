package service

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
	"github.com/alanzhumalin/bank/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	tokenKey       *string
	authRepository repository.AuthRepository
	txManager      repository.TxManagerRepository
	userService    UserService
}

func NewAuthService(tokenKey *string, authRepository repository.AuthRepository, userService UserService, txManager repository.TxManagerRepository) AuthService {
	return &authService{tokenKey: tokenKey, authRepository: authRepository, userService: userService, txManager: txManager}
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

func (a *authService) Login(ctx context.Context, req dto.LoginRequest, ip string, device string) (*dto.TokenPair, error) {

	userDetails, err := a.authRepository.GetDetails(ctx, req.PhoneNumber)

	if err != nil {
		return nil, err
	}
	hashedPassword := userDetails.Password

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return &dto.TokenPair{}, domain.ErrorPasswordNotCorrect
	}
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

func (a *authService) UpdateSession(ctx context.Context, userId int, role string, sessionId string) (*dto.TokenPair, error) {

	accessToken, err := jwt.GeneratateAccessToken(userId, role, sessionId, *a.tokenKey)

	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, _, err := jwt.GeneratateRefreshToken(userId, role, sessionId, *a.tokenKey)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := jwt.HashRefreshToken(refreshToken, *a.tokenKey)

	if err = a.authRepository.Update(ctx, hashedRefreshToken, *expiresAt, sessionId); err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (a *authService) Logout(ctx context.Context, sessionId string) error {
	return a.authRepository.Revoke(ctx, sessionId)
}

func (a *authService) LogoutFromAllDevices(ctx context.Context, userId int) error {
	return a.authRepository.RevokeAllUserDevices(ctx, userId)
}
