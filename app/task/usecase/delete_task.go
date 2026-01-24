package usecase

import (
	"context"
	"project_m_backend/app/task/ports"
	"project_m_backend/apperrors"
)

type DeleteTaskUseCase struct {
	taskRepo ports.TaskRepository
}

func NewDeleteTaskUseCase(taskRepo ports.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (tc *DeleteTaskUseCase) Execute(ctx context.Context, taskId int64, userId int64) *apperrors.AppError {
	existingTask, appErr := tc.taskRepo.GetTaskById(ctx, taskId)
	if appErr != nil {
		return appErr
	}

	if existingTask.AssignedBy != userId {
		return apperrors.NewUnauthorized(nil, "only the manager who assigned this task can delete it")
	}

	appErr = tc.taskRepo.DeleteTask(ctx, taskId)
	if appErr != nil {
		return appErr
	}

	return nil
}
