package usecase

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"github.com/stempo/backend/internal/infrastructure/oauth"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string, name, phone *string, cityID *uint) (*entity.User, string, string, error)
	Login(email, password string) (*entity.User, string, string, error)
	LoginWithGoogle(idToken string) (*entity.User, string, string, error)
	Refresh(refreshToken string) (string, error)
	GetCurrentUser(userID uint) (*entity.User, error)
	UpdateProfile(userID uint, name, phone *string, cityID *uint) (*entity.User, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
	SetPasswordFromInvite(token, password string) (*entity.User, string, string, error)
}

type authUsecase struct {
	userRepo        repository.UserRepository
	googleVerifier  oauth.GoogleTokenVerifier
}

func NewAuthUsecase(userRepo repository.UserRepository, googleVerifier oauth.GoogleTokenVerifier) AuthUsecase {
	return &authUsecase{userRepo: userRepo, googleVerifier: googleVerifier}
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

	if user.Password == "" {
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

func (u *authUsecase) LoginWithGoogle(idToken string) (*entity.User, string, string, error) {
	googleUser, err := u.googleVerifier.Verify(idToken)
	if err != nil {
		return nil, "", "", err
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@stempo.com"
	}
	if strings.EqualFold(googleUser.Email, adminEmail) {
		return nil, "", "", errors.New("this email is reserved for admin")
	}

	user, err := u.userRepo.FindByGoogleID(googleUser.GoogleID)
	if err != nil || user == nil {
		existingUser, emailErr := u.userRepo.FindByEmail(googleUser.Email)
		if emailErr == nil && existingUser != nil {
			user = existingUser
			user.GoogleID = &googleUser.GoogleID
			if user.Name == nil && googleUser.Name != "" {
				name := googleUser.Name
				user.Name = &name
			}
			if user.AvatarURL == nil && googleUser.Picture != "" {
				picture := googleUser.Picture
				user.AvatarURL = &picture
			}
			if err := u.userRepo.Update(user); err != nil {
				return nil, "", "", err
			}
		} else {
			var name *string
			if googleUser.Name != "" {
				n := googleUser.Name
				name = &n
			}
			var avatarURL *string
			if googleUser.Picture != "" {
				p := googleUser.Picture
				avatarURL = &p
			}
			googleID := googleUser.GoogleID
			user = &entity.User{
				Email:     googleUser.Email,
				GoogleID:  &googleID,
				Name:      name,
				AvatarURL: avatarURL,
			}
			if err := u.userRepo.Create(user); err != nil {
				return nil, "", "", err
			}
		}
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

func (u *authUsecase) SetPasswordFromInvite(token, password string) (*entity.User, string, string, error) {
	user, err := u.userRepo.FindByInviteToken(token)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid or expired invite token")
	}

	if user.InviteTokenExpiry != nil && time.Now().After(*user.InviteTokenExpiry) {
		return nil, "", "", errors.New("invite token has expired")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", "", err
	}

	user.Password = hashedPassword
	user.InviteToken = nil
	user.InviteTokenExpiry = nil

	if err := u.userRepo.Update(user); err != nil {
		return nil, "", "", err
	}

	accessToken, err := u.generateToken(user.ID, accessTokenTTL)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := u.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
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

func (u *authUsecase) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	if currentPassword == newPassword {
		return errors.New("new password must be different from current password")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return u.userRepo.Update(user)
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
