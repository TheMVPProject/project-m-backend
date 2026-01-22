package jwt

import (
	"fmt"
	"project_m_backend/domain/user"
	"project_m_backend/pkg/crypto"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	encryptor crypto.Encryptor
}

func NewJWTService(secretKey string, encryptor crypto.Encryptor) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		encryptor: encryptor,
	}
}

func (s *JWTService) GenerateAccessToken(user *user.User, expiry time.Duration) (string, error) {
	encryptedUserID, err := s.encryptor.Encrypt(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt user ID: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  encryptedUserID,
		"type": "access",
		"exp":  time.Now().Add(expiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) GenerateRefreshToken(user *user.User, expiry time.Duration) (string, error) {
	encryptedUserID, err := s.encryptor.Encrypt(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt user ID: %w", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  encryptedUserID,
		"type": "refresh",
		"exp":  time.Now().Add(expiry).Unix(),
	})

	tokenString, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ParseAccessToken(tokenString string) (int64, error) {
	return s.parseToken(tokenString, "access")
}

func (s *JWTService) ParseRefreshToken(tokenString string) (int64, error) {
	return s.parseToken(tokenString, "refresh")
}

func (s *JWTService) parseToken(tokenString string, expectedType string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != expectedType {
			return 0, fmt.Errorf("invalid token type")
		}

		encryptedUserID, ok := claims["sub"].(string)
		if !ok {
			return 0, fmt.Errorf("invalid token claims: sub not found or not a string")
		}

		decryptedUserID, err := s.encryptor.Decrypt(encryptedUserID)
		if err != nil {
			return 0, fmt.Errorf("failed to decrypt user ID: %w", err)
		}

		userID, err := strconv.ParseInt(decryptedUserID, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse decrypted user ID: %w", err)
		}

		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

// GetSecretKey returns the secret key used for signing tokens.
func (s *JWTService) GetSecretKey() []byte {
	return []byte(s.secretKey)
}

// DecryptUserID decrypts the encrypted user ID string and returns the int64 representation.
func (s *JWTService) DecryptUserID(encryptedUserID string) (int64, error) {
	decryptedUserID, err := s.encryptor.Decrypt(encryptedUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt user ID: %w", err)
	}

	userID, err := strconv.ParseInt(decryptedUserID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse decrypted user ID: %w", err)
	}

	return userID, nil
}
