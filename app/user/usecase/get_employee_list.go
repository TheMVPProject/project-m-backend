package usecase

import (
	"context"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/app/user/ports"
	"project_m_backend/apperrors"
)

type GetEmployeesListUseCase struct{
	userRepo ports.UserRepository
}

func NewGetEmployeeListUseCase(userRepo ports.UserRepository) *GetEmployeesListUseCase{
	return &GetEmployeesListUseCase{userRepo: userRepo,}
}

func (uc *GetEmployeesListUseCase)Execute(ctx context.Context) ([]*appModel.AppUser, *apperrors.AppError){
	employees, apperr := uc.userRepo.GetEmployeesList(ctx)
	if apperr != nil{
		return nil, apperr
	}
	return  employees, nil
}