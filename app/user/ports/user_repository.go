package ports

import (
	"context"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/apperrors"
)

type UserRepository interface{
	GetUserById(ctx context.Context, id int64) (*appModel.AppUser, *apperrors.AppError)
	GetUserByEmail(ctx context.Context, email string) (*appModel.AppUser, *apperrors.AppError)
	CreateUser(ctx context.Context, appUser *appModel.AppUser) (int64, *apperrors.AppError)
}