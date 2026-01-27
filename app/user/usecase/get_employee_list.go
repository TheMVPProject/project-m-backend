package usecase

import (
	"context"
	appModel "assignly/app/user/model"
	"assignly/app/user/ports"
	"assignly/apperrors"
)

type GetEmployeesListUseCase struct{
	userRepo ports.UserRepository
}

func NewGetEmployeeListUseCase(userRepo ports.UserRepository) *GetEmployeesListUseCase{
	return &GetEmployeesListUseCase{userRepo: userRepo,}
}

func (uc *GetEmployeesListUseCase)Execute(ctx context.Context) ([]*appModel.EmployeeUser, *apperrors.AppError){
	employees, apperr := uc.userRepo.GetEmployeesList(ctx)
	if apperr != nil{
		return nil, apperr
	}
	return  employees, nil
}