package model

import (
	"time"
	"project_m_backend/domain/user"
)

type AppUser struct{
	ID int64
	FirstName string
	LastName string
	Email string
	PasswordHash string
	Position user.UserType
	CreatedAt time.Time
	LastSignIn time.Time
}