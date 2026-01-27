package ports

import (
	"context"
	appModel "assignly/app/user/model"
	"assignly/apperrors"
)

type UserRepository interface{
	GetUserById(ctx context.Context, id int64) (*appModel.AppUser, *apperrors.AppError)
	GetUserByEmail(ctx context.Context, email string) (*appModel.AppUser, *apperrors.AppError)
	CreateUser(ctx context.Context, appUser *appModel.AppUser) *apperrors.AppError
	GetEmployeesList(ctx context.Context) ([]*appModel.EmployeeUser, *apperrors.AppError)
}