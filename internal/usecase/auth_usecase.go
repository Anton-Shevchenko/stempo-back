package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string, name, phone *string, cityID *uint) (*entity.User, string, string, error)
	Login(email, password string) (*entity.User, string, string, error)
	Refresh(refreshToken string) (string, error)
	GetCurrentUser(userID uint) (*entity.User, error)
	UpdateProfile(userID uint, name, phone *string, cityID *uint) (*entity.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

const (
	accessTokenTTL  = 24 * time.Hour      // 1 day
	refreshTokenTTL = 365 * 24 * time.Hour // 1 year
)

func (u *authUsecase) Register(email, password string, name, phone *string, cityID *uint) (*entity.User, string, string, error) {
	existingUser, err := u.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, "", "", errors.New("user with this email already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", "", err
	}

	user := &entity.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Phone:    phone,
		CityID:   cityID,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, "", "", err
	}

	token, err := u.generateToken(user.ID, accessTokenTTL)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := u.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, token, refreshToken, nil
}

func (u *authUsecase) Login(email, password string) (*entity.User, string, string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	token, err := u.generateToken(user.ID, accessTokenTTL)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := u.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, token, refreshToken, nil
}

func (u *authUsecase) Refresh(refreshToken string) (string, error) {
	userID, err := u.parseRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	token, err := u.generateToken(userID, accessTokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUsecase) GetCurrentUser(userID uint) (*entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *authUsecase) UpdateProfile(userID uint, name, phone *string, cityID *uint) (*entity.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if name != nil {
		user.Name = name
	}
	if phone != nil {
		user.Phone = phone
	}
	if cityID != nil {
		user.CityID = cityID
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) generateToken(userID uint, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *authUsecase) generateRefreshToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(refreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *authUsecase) parseRefreshToken(tokenString string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	if claims["type"] != "refresh" {
		return 0, errors.New("not a refresh token")
	}

	var userID uint
	if userIDFloat, ok := claims["user_id"].(float64); ok {
		userID = uint(userIDFloat)
	} else if userIDInt, ok := claims["user_id"].(int); ok {
		userID = uint(userIDInt)
	} else {
		return 0, errors.New("invalid user id in token")
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
