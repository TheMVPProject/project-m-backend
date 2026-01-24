package usecase

import (
	"context"
	appModel "assignly/app/user/model"
	"assignly/app/user/ports"
	"assignly/apperrors"
	domain "assignly/domain/user"
)


type LoginWithEmailUseCase struct {
	userRepo    ports.UserRepository
	hasher      domain.PasswordHasher
	userFactory *domain.UserFactory
}


func NewLoginWithEmailUseCase(userRepo ports.UserRepository, hasher domain.PasswordHasher,
userFactory *domain.UserFactory) *LoginWithEmailUseCase{
	return &LoginWithEmailUseCase{
		userRepo: userRepo,
		hasher: hasher,
		userFactory: userFactory,
	}
}


func (uc *LoginWithEmailUseCase) Execute(ctx context.Context, email, password string) (*appModel.AppUser, *apperrors.AppError) {
	if email == ""{
		return nil, apperrors.NewEmptyEmail(nil)
	}

	if password == ""{
		return nil, apperrors.NewEmptyPassword(nil)
	}

	tempUser, err := uc.userFactory.CreateUnsafeUser(0, "temp", "", email, "", 0)
	if err != nil{
		return  nil, apperrors.NewInvalidEmailFormat(err)
	}
	if err := tempUser.ValidateEmailFormat(uc.userFactory.Validator); err != nil {
		return nil, apperrors.NewInvalidEmailFormat(err)
	}

	appUser, appErr := uc.userRepo.GetUserByEmail(ctx, email)
	if appErr != nil{
		return nil, appErr
	}

	domainUser := appModel.ToDomainUser(appUser)

	err = domainUser.CheckPassword(password, uc.hasher)
	if err != nil{
		return nil, apperrors.NewIncorrectPassword(err)
	}

	return appUser, nil
}