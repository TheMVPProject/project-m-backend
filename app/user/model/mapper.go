package model

import (
	"project_m_backend/domain/user"
	"time"
)

func ToDomainUser(appUser *AppUser) *user.User{
	return &user.User{
		ID: appUser.ID,
		FirstName: appUser.FirstName,
		LastName: appUser.LastName,
		Email: appUser.Email,
		PasswordHash: appUser.PasswordHash,
		Position: appUser.Position,
	}
}

func FromDomainUser(domainUser *user.User) *AppUser{
	return &AppUser{
		ID: domainUser.ID,
		FirstName: domainUser.FirstName,
		LastName: domainUser.LastName,
		Email: domainUser.Email,
		PasswordHash: domainUser.PasswordHash,
		Position: domainUser.Position,
		CreatedAt: time.Now(),
	}
}