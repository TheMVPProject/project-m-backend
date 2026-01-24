package ports

import (
	"context"
	appModel "assignly/app/task/model"
	"assignly/apperrors"
)

type TaskRepository interface{
	AssignTask(ctx context.Context,appTask *appModel.AppTask) *apperrors.AppError
	GetAssignedTaskByUserId(ctx context.Context, user_id int64, isByManegerId bool)([]*appModel.AppTask, *apperrors.AppError)
	GetTaskById(ctx context.Context, taskId int64) (*appModel.AppTask, *apperrors.AppError)
	UpdateTaskStatus(ctx context.Context, taskId int64, status string) *apperrors.AppError
	DeleteTask(ctx context.Context, taskId int64) *apperrors.AppError
}