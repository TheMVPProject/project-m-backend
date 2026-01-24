package model

import (
	"project_m_backend/domain/task"
	"time"
)

type AppTask struct{
	ID int64
	Title string
	Description string
	AssignedBy int64
	AssignedTo int64
	AssignedByName string
	AssignedToName string
	Status task.TaskStatus
	Priority task.TaskPriority
	DueDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}