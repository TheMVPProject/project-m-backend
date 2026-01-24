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
	// Check if user with this email already exists
	existingUser, appErr := uc.userRepo.GetUserByEmail(ctx, email)
	if appErr != nil {
		if appErr.Code != apperrors.CodeNotFound {
			return appErr
		}
	}
	if existingUser != nil {
		return apperrors.NewConflict(nil, "User with this email already exists")
	}

	newUser, err := uc.factory.CreateUser(0, firstName, lastName, email, password, position)
	if err != nil{
		return apperrors.NewInternal(err, "faild to create user entity")
	}

	appUserToCreate := appModel.FromDomainUser(newUser)

	appError := uc.userRepo.CreateUser(ctx, appUserToCreate)
	
	if appError != nil{
		return appError
	}
	return nil
}