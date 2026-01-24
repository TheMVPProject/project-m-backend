package usecase

import (
	"context"
	appModel "assignly/app/user/model"
	"assignly/app/user/ports"
	"assignly/apperrors"
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

