package task

import (
	"time"
)

type TaskStatus string
type TaskPriority int

const(
	TaskStatusPending TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted TaskStatus = "COMPLETED"

	TaskPriorityHigh TaskPriority = 1
	TaskPriorityMedium TaskPriority = 2
	TaskPriorityLow TaskPriority = 3
)

type Task struct{
	ID int64 	`json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	AssignedBy int64 `json:"assigned_by"`
	AssignedByName string `json:"assigned_by_name"`
	AssignedTo int64 `json:"assigned_to"`
	AssignedToName string `json:"assigned_to_name"`
	Status TaskStatus `json:"status"`
	Priority TaskPriority `json:"priority"`
	DueDate time.Time `json:"due_date"`
}

func NewTask(id int64, title, description string, assignedBy int64, assignedByName string, assignedTo int64, assignedToName string, status TaskStatus, priority TaskPriority, dueDate time.Time) (*Task, error){
	task := &Task{
		ID: id,
		Title: title,
		Description: description,
		AssignedBy: assignedBy,
		AssignedByName: assignedByName,
		AssignedTo: assignedTo,
		AssignedToName: assignedToName,
		Status: status,
		Priority: priority,
		DueDate: dueDate,
	}
	return task, nil
}