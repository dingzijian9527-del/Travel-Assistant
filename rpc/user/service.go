package rpcuser

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	user "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/user"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/jwtx"
)

type userService struct {
	repo userRepository
	auth config.AuthConfig
}

func newUserService(repo userRepository) *userService {
	return newUserServiceWithAuth(repo, config.AuthConfig{JWTSecret: "test-jwt-secret", JWTExpire: time.Hour})
}

func newUserServiceWithAuth(repo userRepository, auth config.AuthConfig) *userService {
	if auth.JWTSecret == "" {
		auth.JWTSecret = "travel-assistant-dev-secret"
	}
	if auth.JWTExpire <= 0 {
		auth.JWTExpire = 24 * time.Hour
	}
	return &userService{repo: repo, auth: auth}
}

func (l *userService) Register(ctx context.Context, req *user.RegisterRequest) (*userModel, *serviceError) {
	phone := normalizePhone(req.Phone)
	if !validPhone(phone) {
		return nil, errParam("valid phone is required")
	}
	if len(req.Password) < 8 {
		return nil, errParam("password must be at least 8 characters")
	}
	if _, exists := l.repo.FindByPhone(ctx, phone); exists {
		return nil, errBiz("phone already registered")
	}
	now := time.Now()
	createdUser := &userModel{
		Phone:         phone,
		PasswordHash:  hashPassword(req.Password),
		Nickname:      defaultNickname(phone, req.Nickname),
		AvatarURL:     optionalTrimmed(req.AvatarUrl),
		HomeCity:      optionalTrimmed(req.HomeCity),
		CurrentCity:   optionalTrimmed(req.CurrentCity),
		MemberLevel:   "normal",
		AccountStatus: 1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	storedUser, err := l.repo.Create(ctx, createdUser)
	if err != nil {
		return nil, errBiz("register user failed")
	}
	return storedUser, nil
}

func (l *userService) Login(ctx context.Context, req *user.LoginRequest) (*userModel, string, *serviceError) {
	phone := normalizePhone(req.Phone)
	storedUser, exists := l.repo.FindByPhone(ctx, phone)
	if !exists || !verifyPassword(storedUser.PasswordHash, req.Password) {
		return nil, "", errAuth("phone or password is invalid")
	}
	token, err := jwtx.Generate(jwtx.Config{Secret: l.auth.JWTSecret, Expire: l.auth.JWTExpire}, jwtx.Claims{UserID: storedUser.ID, Phone: storedUser.Phone})
	if err != nil {
		return nil, "", errBiz("generate login token failed")
	}
	return storedUser, token, nil
}

func (l *userService) GetProfile(ctx context.Context, userID string) (*userModel, *serviceError) {
	if strings.TrimSpace(userID) == "" {
		return nil, errParam("valid id is required")
	}
	storedUser, exists := l.repo.FindByID(ctx, userID)
	if !exists {
		return nil, errNotFound("user not found")
	}
	return storedUser, nil
}

func (l *userService) UpdateProfile(ctx context.Context, req *user.UpdateProfileRequest) (*userModel, *serviceError) {
	if strings.TrimSpace(req.Id) == "" {
		return nil, errParam("valid id is required")
	}
	updatedUser, exists := l.repo.Update(ctx, req.Id, func(storedUser *userModel) {
		if req.Nickname != nil {
			storedUser.Nickname = strings.TrimSpace(*req.Nickname)
		}
		if req.AvatarUrl != nil {
			storedUser.AvatarURL = strings.TrimSpace(*req.AvatarUrl)
		}
		if req.HomeCity != nil {
			storedUser.HomeCity = strings.TrimSpace(*req.HomeCity)
		}
		if req.CurrentCity != nil {
			storedUser.CurrentCity = strings.TrimSpace(*req.CurrentCity)
		}
		storedUser.UpdatedAt = time.Now()
	})
	if !exists {
		return nil, errNotFound("user not found")
	}
	return updatedUser, nil
}

func hashPassword(password string) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		sum := sha256.Sum256([]byte("travel-assistant:" + password))
		return "legacy$" + hex.EncodeToString(sum[:])
	}
	hash := stretchPassword(password, salt)
	return "v1$" + base64.RawURLEncoding.EncodeToString(salt) + "$" + base64.RawURLEncoding.EncodeToString(hash)
}

func verifyPassword(storedHash string, password string) bool {
	parts := strings.Split(storedHash, "$")
	if len(parts) == 3 && parts[0] == "v1" {
		salt, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return false
		}
		expected, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			return false
		}
		actual := stretchPassword(password, salt)
		return subtle.ConstantTimeCompare(actual, expected) == 1
	}
	if strings.HasPrefix(storedHash, "legacy$") {
		sum := sha256.Sum256([]byte("travel-assistant:" + password))
		return storedHash == "legacy$"+hex.EncodeToString(sum[:])
	}
	return false
}

func stretchPassword(password string, salt []byte) []byte {
	hash := sha256.Sum256(append(append([]byte{}, salt...), []byte(password)...))
	out := hash[:]
	for i := 0; i < 12000; i++ {
		next := sha256.Sum256(append(out, []byte(password)...))
		out = next[:]
	}
	return out
}
