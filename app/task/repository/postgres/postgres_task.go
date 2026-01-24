package postgres

import (
	"context"
	"database/sql"
	"fmt"
	appModel "project_m_backend/app/task/model"
	"project_m_backend/apperrors"
	"time"
)

const (
	queryAssignTask = `INSERT INTO tasks
	(title, description, assigned_by, assigned_to,assigned_by_name,assigned_to_name, status, priority, due_date,
	created_at, updated_at)VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`
	queryGetTasksByEmployeeID = `SELECT id, title, description, assigned_by, assigned_to,assigned_by_name,assigned_to_name, status, priority, due_date, created_at, updated_at FROM tasks WHERE assigned_to = $1;`
	queryGetTasksByManegerID  = `SELECT id, title, description, assigned_by,assigned_to, assigned_by_name,assigned_to_name, status, priority, due_date, created_at, updated_at FROM tasks WHERE assigned_by = $1;`
	queryGetTaskByID          = `SELECT id, title, description, assigned_by, assigned_to,assigned_by_name,assigned_to_name, status, priority, due_date, created_at, updated_at FROM tasks WHERE id = $1;`
	queryUpdateTaskStatus     = `UPDATE tasks SET status = $1, updated_at = $2 WHERE id = $3;`
	queryDeleteTask           = `DELETE FROM tasks WHERE id = $1;`
)

type PostgresTaskRepository struct {
	db                      *sql.DB
	assignTaskStmt          *sql.Stmt
	getTaskByEmployeeIdStmt *sql.Stmt
	getTaskByManegerIdStmt  *sql.Stmt
	getTaskByIdStmt         *sql.Stmt
	updateTaskStatusStmt    *sql.Stmt
	deleteTaskStmt          *sql.Stmt
}

func NewPostgresTaskRepository(db *sql.DB) (*PostgresTaskRepository, error) {
	createAssignTaskStmt, err := db.PrepareContext(context.Background(), queryAssignTask)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error Preparing assigne task")
	}

	createGetTasksByEmployeeIDStmt, err := db.PrepareContext(context.Background(), queryGetTasksByEmployeeID)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error preparing get tasks by Employee id")
	}

	createGetTasksByManegerIDStmt, err := db.PrepareContext(context.Background(), queryGetTasksByManegerID)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error preparing get tasks by Maneger id")
	}

	createGetTaskByIDStmt, err := db.PrepareContext(context.Background(), queryGetTaskByID)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error preparing get task by id")
	}

	createUpdateTaskStatusStmt, err := db.PrepareContext(context.Background(), queryUpdateTaskStatus)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error preparing update task status")
	}

	createDeleteTaskStmt, err := db.PrepareContext(context.Background(), queryDeleteTask)
	if err != nil {
		return nil, apperrors.NewInternal(err, "error preparing delete task")
	}

	return &PostgresTaskRepository{
		db:                      db,
		assignTaskStmt:          createAssignTaskStmt,
		getTaskByEmployeeIdStmt: createGetTasksByEmployeeIDStmt,
		getTaskByManegerIdStmt:  createGetTasksByManegerIDStmt,
		getTaskByIdStmt:         createGetTaskByIDStmt,
		updateTaskStatusStmt:    createUpdateTaskStatusStmt,
		deleteTaskStmt:          createDeleteTaskStmt,
	}, nil
}

func (r *PostgresTaskRepository) Close() error {
	var errs []error
	if err := r.assignTaskStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := r.getTaskByEmployeeIdStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := r.getTaskByManegerIdStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := r.getTaskByIdStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := r.updateTaskStatusStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := r.deleteTaskStmt.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return apperrors.NewInternal(fmt.Errorf("Multiple errors: %v", errs),
			"errors closing statements")
	}

	return nil
}

func (r *PostgresTaskRepository) AssignTask(ctx context.Context, task *appModel.AppTask) *apperrors.AppError {
	createdAt := time.Now()

	var id int64
	err := r.assignTaskStmt.QueryRowContext(
		ctx, task.Title, task.Description, task.AssignedBy, task.AssignedTo, task.AssignedByName, task.AssignedToName,
		task.Status, task.Priority, task.DueDate, createdAt, createdAt,
	).Scan(&id)

	if err != nil {
		return apperrors.NewInternal(err, "faild to Assign a task")
	}
	return nil
}

func (r *PostgresTaskRepository) GetAssignedTaskByUserId(ctx context.Context, userId int64, isByManegerId bool) ([]*appModel.AppTask, *apperrors.AppError) {
	var (
		rows *sql.Rows
		err  error
	)

	if !isByManegerId {
		rows, err = r.getTaskByEmployeeIdStmt.QueryContext(ctx, userId)
	} else {
		rows, err = r.getTaskByManegerIdStmt.QueryContext(ctx, userId)
	}

	if err != nil {
		return nil, apperrors.NewInternal(err, "failed to query assigned task")
	}
	defer rows.Close()

	var tasks []*appModel.AppTask

	for rows.Next() {
		task := &appModel.AppTask{}

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.AssignedBy,
			&task.AssignedTo,
			&task.AssignedByName,
			&task.AssignedToName,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err.Error())
			return nil, apperrors.NewInternal(
				err, "faild to scan tasks")
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.NewInternal(
			err, "row iteration error")
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) GetTaskById(ctx context.Context, taskId int64) (*appModel.AppTask, *apperrors.AppError) {
	task := &appModel.AppTask{}

	err := r.getTaskByIdStmt.QueryRowContext(ctx, taskId).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.AssignedBy,
		&task.AssignedTo,
		&task.AssignedByName,
		&task.AssignedToName,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewNotFound(err, "task not found")
		}
		return nil, apperrors.NewInternal(err, "failed to get task by id")
	}

	return task, nil
}

func (r *PostgresTaskRepository) UpdateTaskStatus(ctx context.Context, taskId int64, status string) *apperrors.AppError {
	updatedAt := time.Now()

	result, err := r.updateTaskStatusStmt.ExecContext(ctx, status, updatedAt, taskId)
	if err != nil {
		return apperrors.NewInternal(err, "failed to update task status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternal(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return apperrors.NewNotFound(nil, "task not found")
	}

	return nil
}

func (r *PostgresTaskRepository) DeleteTask(ctx context.Context, taskId int64) *apperrors.AppError {
	result, err := r.deleteTaskStmt.ExecContext(ctx, taskId)
	if err != nil {
		return apperrors.NewInternal(err, "failed to delete task")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewInternal(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return apperrors.NewNotFound(nil, "task not found")
	}

	return nil
}
