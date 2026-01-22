package usecase

import (
	"context"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/app/user/ports"
	"project_m_backend/apperrors"
	"project_m_backend/domain/user"
	userDomain "project_m_backend/domain/user"
)


type CreatUserUseCase struct{
	userRepo ports.UserRepository
	factory *userDomain.UserFactory
}

func NewCreateUserUseCase(userRepo ports.UserRepository, factory *userDomain.UserFactory) *CreatUserUseCase{
	return &CreatUserUseCase{
		userRepo: userRepo,
		factory: factory,
	}
}

func(uc *CreatUserUseCase) Execute(ctx context.Context, firstName, lastName, email, password string, position user.UserType) *apperrors.AppError{
	newUser, err := uc.factory.CreateUser(0, firstName, lastName, email, password, position)
	if err != nil{
		return apperrors.NewInternal(err, "faild to create user entity")
	}

	appUserToCreate := appModel.FromDomainUser(newUser)

	_, appErr := uc.userRepo.CreateUser(ctx, appUserToCreate)
	if appErr != nil{
		return appErr
	}
	return nil
}