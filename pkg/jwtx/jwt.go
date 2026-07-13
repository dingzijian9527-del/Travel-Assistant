package jwtx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Config 描述令牌签名和过期时间。
type Config struct {
	Secret string
	Expire time.Duration
}

// Claims 描述旅行助手登录态中的用户身份。
type Claims struct {
	UserID    string
	Phone     string
	Role      string
	Expire    int64
	TokenType string
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwtPayload struct {
	UserID    string `json:"user_id"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty"`
	Expire    int64  `json:"exp"`
	TokenType string `json:"token_type"`
}

// Generate 生成使用 HS256 签名的令牌。
func Generate(cfg Config, claims Claims) (string, error) {
	return generate(cfg, claims, TokenTypeAccess)
}

func GenerateRefresh(cfg Config, claims Claims, expire time.Duration) (string, error) {
	cfg.Expire = expire
	return generate(cfg, claims, TokenTypeRefresh)
}

func GeneratePair(cfg Config, claims Claims, refreshExpire time.Duration) (string, string, error) {
	access, err := Generate(cfg, claims)
	if err != nil {
		return "", "", err
	}
	refresh, err := GenerateRefresh(cfg, claims, refreshExpire)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func generate(cfg Config, claims Claims, tokenType string) (string, error) {
	if strings.TrimSpace(cfg.Secret) == "" {
		return "", errors.New("jwt secret is required")
	}
	if cfg.Expire <= 0 {
		cfg.Expire = 24 * time.Hour
	}
	claims.Expire = time.Now().Add(cfg.Expire).Unix()
	headerPart, err := encodeJSON(jwtHeader{Algorithm: "HS256", Type: "JWT"})
	if err != nil {
		return "", err
	}
	payloadPart, err := encodeJSON(jwtPayload{UserID: claims.UserID, Phone: claims.Phone, Role: claims.Role, Expire: claims.Expire, TokenType: tokenType})
	if err != nil {
		return "", err
	}
	unsigned := headerPart + "." + payloadPart
	return unsigned + "." + sign(unsigned, cfg.Secret), nil
}

// Parse 校验令牌签名和过期时间，并返回用户身份。
func Parse(cfg Config, token string) (Claims, error) {
	return parseWithType(cfg, token, TokenTypeAccess)
}

func ParseRefresh(cfg Config, token string) (Claims, error) {
	return parseWithType(cfg, token, TokenTypeRefresh)
}

func parseWithType(cfg Config, token string, expectedType string) (Claims, error) {
	if strings.TrimSpace(cfg.Secret) == "" {
		return Claims{}, errors.New("jwt secret is required")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errors.New("invalid jwt format")
	}
	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(unsigned, cfg.Secret)), []byte(parts[2])) {
		return Claims{}, errors.New("invalid jwt signature")
	}
	var payload jwtPayload
	if err := decodeJSON(parts[1], &payload); err != nil {
		return Claims{}, err
	}
	if payload.UserID == "" {
		return Claims{}, errors.New("jwt user id is required")
	}
	if payload.TokenType != expectedType {
		return Claims{}, errors.New("jwt token type is invalid")
	}
	if payload.Expire <= time.Now().Unix() {
		return Claims{}, errors.New("jwt expired")
	}
	return Claims{UserID: payload.UserID, Phone: payload.Phone, Role: payload.Role, Expire: payload.Expire, TokenType: payload.TokenType}, nil
}

func encodeJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeJSON(part string, value any) error {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return fmt.Errorf("decode jwt part failed: %w", err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("parse jwt payload failed: %w", err)
	}
	return nil
}

func sign(unsigned string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// UserIDInt64 将字符串用户编号转换为整数，供智能体接口复用。
func UserIDInt64(userID string) (int64, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(userID), "user-")
	if value, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return value, nil
	}
	if trimmed == "" {
		return 0, errors.New("user id is required")
	}
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(trimmed))
	return int64(hash.Sum64() & 0x7fffffffffffffff), nil
}
