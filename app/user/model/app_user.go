package model

import (
	"time"
	"assignly/domain/user"
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

type EmployeeUser struct{
	ID int64
	Name string
}