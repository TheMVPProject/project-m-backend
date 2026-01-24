package model

import (
	"assignly/domain/task"
	"time"
)

func ToDomainUser(appTask *AppTask) *task.Task{
	return &task.Task{
		ID: appTask.ID,
		Title: appTask.Title,
		Description: appTask.Description,
		AssignedBy: appTask.AssignedBy,
		AssignedTo: appTask.AssignedTo,
		AssignedByName: appTask.AssignedByName,
		AssignedToName: appTask.AssignedToName,
		Status: appTask.Status,
		Priority: appTask.Priority,
		DueDate: appTask.DueDate,
	}
}

func FromDomainUser(task *task.Task) *AppTask{
	return  &AppTask{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		AssignedBy: task.AssignedBy,
		AssignedByName: task.AssignedByName,
		AssignedTo: task.AssignedTo,
		AssignedToName: task.AssignedToName,
		Status: task.Status,
		Priority: task.Priority,
		DueDate: task.DueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}