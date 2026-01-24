package ports

import (
	"context"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/apperrors"
)

type UserRepository interface{
	GetUserById(ctx context.Context, id int64) (*appModel.AppUser, *apperrors.AppError)
	GetUserByEmail(ctx context.Context, email string) (*appModel.AppUser, *apperrors.AppError)
	CreateUser(ctx context.Context, appUser *appModel.AppUser) *apperrors.AppError
	GetEmployeesList(ctx context.Context) ([]*appModel.AppUser, *apperrors.AppError)
}