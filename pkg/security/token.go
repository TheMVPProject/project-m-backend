package security

import (
	"fmt"
	"math/rand"
	"project_m_backend/apperrors"
	"strconv"
	"time"
)


const digits = "0123456789"

type TokenGenerator struct{}

func NewTokenGenerator() *TokenGenerator{
	return &TokenGenerator{}
}

func (tg *TokenGenerator) GenerateSessionToken(userId int64) string{
	unixTimeString := fmt.Sprintf("%013d", time.Now().UnixMilli())
	userIdStr := strconv.FormatInt(userId, 10)
	sessionTokenStr := userIdStr + unixTimeString
	return sessionTokenStr
}

func (tg *TokenGenerator) GenerateOTP(length int) (string, *apperrors.AppError){
	b := make([]byte, length)
	for i := range b{
		n := rand.Int63n(9)
		b[i] = digits[n]
	}
	return string(b), nil
}