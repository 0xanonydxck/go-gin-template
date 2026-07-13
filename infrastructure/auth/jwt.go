package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chai-rs/simple-bookstore/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager defines methods for JWT token operations.
type TokenManager interface {
	CreateToken(userId, email string) (*TokenProperties, error)
	ExtractTokenMetadata(*http.Request) (*AccessProperties, error)
}

type tokenManager struct{}

// NewTokenManager creates a new TokenManager instance.
func NewTokenManager() TokenManager {
	return &tokenManager{}
}

// CreateToken generates new access and refresh tokens for a user.
func (t *tokenManager) CreateToken(userId, email string) (*TokenProperties, error) {
	properties := new(TokenProperties)

	now := time.Now()
	properties.AccessTokenExpire = now.Add(time.Minute * 30).Unix()
	properties.AccessTokenUUID = uuid.New().String()
	properties.RefreshTokenExpire = now.Add(time.Hour * 24 * 7).Unix()
	properties.RefreshTokenUUID = ToRefreshUUID(properties.RefreshTokenUUID, userId)

	// Create access token
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = properties.AccessTokenUUID
	atClaims["user_id"] = userId
	atClaims["email"] = email
	atClaims["exp"] = properties.AccessTokenExpire
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	properties.AccessToken, err = at.SignedString([]byte(config.ACCESS_SECRET))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = properties.RefreshTokenUUID
	rtClaims["user_id"] = userId
	rtClaims["email"] = email
	rtClaims["exp"] = properties.RefreshTokenExpire
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	properties.RefreshToken, err = rt.SignedString([]byte(config.REFRESH_SECRET))
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// ExtractTokenMetadata extracts token metadata from HTTP request.
func (t *tokenManager) ExtractTokenMetadata(r *http.Request) (*AccessProperties, error) {
	return ExtractTokenMetadata(r)
}

// TokenValid checks if the token in the request is valid.
func TokenValid(r *http.Request) error {
	token, err := VerifyAuthorizationHeader(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// VerifyAuthorizationHeader verifies and parses the Authorization header.
func VerifyAuthorizationHeader(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	return VerifyToken(tokenString)
}

// VerifyToken verifies and parses a JWT token string.
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return []byte(config.ACCESS_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// ExtractToken retrieves the JWT token from the Authorization header.
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

// Extract retrieves access properties from a JWT token.
func Extract(token *jwt.Token) (*AccessProperties, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid access uuid")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email")
	}

	properties := &AccessProperties{
		TokenUUID: accessUUID,
		UserID:    userID,
		Email:     email,
	}

	return properties, nil
}

// ExtractTokenMetadata extracts token metadata from HTTP request.
// This is a package-level function that implements the TokenManager interface.
func ExtractTokenMetadata(r *http.Request) (*AccessProperties, error) {
	token, err := VerifyAuthorizationHeader(r)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	acc, err := Extract(token)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
