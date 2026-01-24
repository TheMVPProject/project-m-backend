package usecase

import (
	"context"
	"project_m_backend/app/task/ports"
	"project_m_backend/apperrors"
	"project_m_backend/domain/task"
	"time"
	appModel "project_m_backend/app/task/model"
)

type AssigneTaskUseCase struct{
	taskRepo ports.TaskRepository
}

func NewAssignTaskUseCase(taskRepo ports.TaskRepository) *AssigneTaskUseCase{
	return &AssigneTaskUseCase{
		taskRepo: taskRepo,
	}
}

func(tc *AssigneTaskUseCase) Execute(ctx context.Context, title, description, assignedByName, assignedToName string, assignedBy, assignedTo int64, status task.TaskStatus, priority task.TaskPriority, dueDate time.Time) *apperrors.AppError{
	newTask, err := task.NewTask(0, title, description, assignedBy, assignedByName, assignedTo, assignedToName, status, priority, dueDate)
	if err != nil{
		return apperrors.NewInternal(err, "faild to create task entity")
	}

	appTaskToCreate := appModel.FromDomainUser(newTask)

	appErr := tc.taskRepo.AssignTask(ctx, appTaskToCreate)
	if appErr != nil{
		return appErr
	}

	return nil
}