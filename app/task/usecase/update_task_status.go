package usecase

import (
	"context"
	"assignly/app/task/ports"
	"assignly/apperrors"
	"assignly/domain/task"
)

type UpdateTaskStatusUseCase struct {
	taskRepo ports.TaskRepository
}

func NewUpdateTaskStatusUseCase(taskRepo ports.TaskRepository) *UpdateTaskStatusUseCase {
	return &UpdateTaskStatusUseCase{
		taskRepo: taskRepo,
	}
}

func (tc *UpdateTaskStatusUseCase) Execute(ctx context.Context, taskId int64, userId int64, status task.TaskStatus) *apperrors.AppError {
	existingTask, appErr := tc.taskRepo.GetTaskById(ctx, taskId)
	if appErr != nil {
		return appErr
	}

	if existingTask.AssignedTo != userId && existingTask.AssignedBy != userId {
		return apperrors.NewUnauthorized(nil, "you are not authorized to update this task status")
	}
	println("now updating task status")
	appErr = tc.taskRepo.UpdateTaskStatus(ctx, taskId, string(status))
	if appErr != nil {
		return appErr
	}

	return nil
}
