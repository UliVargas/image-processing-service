package auth

import (
	"image-processing-service/internal/shared/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = utils.NewError(401, "INVALID_TOKEN", "Token invalido", nil)
)

type AppClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
}

func NewTokenManager(secret string, expiry time.Duration) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secret),
		expiry:    expiry,
		issuer:    "image-processing-service",
	}
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	JTI          string
}

func (m *TokenManager) GeneratePair(UserID string) (*TokenResponse, error) {
	jti := utils.GenerateID()
	refreshToken := utils.GenerateID()

	claims := AppClaims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		JTI:          jti,
	}, nil
}

func (m *TokenManager) Validate(tokenStr string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AppClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
