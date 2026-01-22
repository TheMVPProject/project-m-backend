package usecase

import (
	"context"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/app/user/ports"
	"project_m_backend/apperrors"
)


type GetUserByIdUseCase struct{
	userRepo ports.UserRepository
}

func NewGetUserByIDUseCase(userRepo ports.UserRepository) *GetUserByIdUseCase{
	return  &GetUserByIdUseCase{userRepo: userRepo}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, id int64) (*appModel.AppUser, *apperrors.AppError){
	appUser, appErr := uc.userRepo.GetUserById(ctx, id)
	if appErr != nil{
		return nil, appErr
	}
	return appUser, nil
}

