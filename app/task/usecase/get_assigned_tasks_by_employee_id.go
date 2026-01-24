package usecase

import (
	"context"
	appModel "project_m_backend/app/task/model"
	"project_m_backend/app/task/ports"
	"project_m_backend/apperrors"
)

type GetAssignedTaskByEmployeeIdUseCase struct{
	taskRepo ports.TaskRepository
}

func NewGetAssignedTaskByEmployeeIdUseCase (taskRepo ports.TaskRepository) *GetAssignedTaskByEmployeeIdUseCase{
	return  &GetAssignedTaskByEmployeeIdUseCase{taskRepo: taskRepo}
}

func (tc *GetAssignedTaskByEmployeeIdUseCase) Execute(ctx context.Context, userId int64, isByManegerId bool) ([]*appModel.AppTask, *apperrors.AppError){
	tasks, apperr := tc.taskRepo.GetAssignedTaskByUserId(ctx, userId, isByManegerId)
	if apperr != nil{
		return nil, apperr
	}
	return tasks, nil
}