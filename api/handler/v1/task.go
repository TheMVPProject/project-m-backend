package v1

import (
	"strconv"

	"assignly/app/task/usecase"
	userUseCase "assignly/app/user/usecase"
	"assignly/domain/task"
	"assignly/domain/user"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct{
	assignTask *usecase.AssigneTaskUseCase
	getAssignedTaskByEmployeeId *usecase.GetAssignedTaskByEmployeeIdUseCase
	getUserByIdUseCase *userUseCase.GetUserByIdUseCase
	updateTaskStatus *usecase.UpdateTaskStatusUseCase
	deleteTask *usecase.DeleteTaskUseCase
}


func NewTaskHandler(assignTask *usecase.AssigneTaskUseCase,
	getAssignedTaskByEmployeeId *usecase.GetAssignedTaskByEmployeeIdUseCase,
	getUserByIdUseCase *userUseCase.GetUserByIdUseCase,
	updateTaskStatus *usecase.UpdateTaskStatusUseCase,
	deleteTask *usecase.DeleteTaskUseCase) *TaskHandler{
		return &TaskHandler{
			assignTask: assignTask,
			getAssignedTaskByEmployeeId: getAssignedTaskByEmployeeId,
			getUserByIdUseCase: getUserByIdUseCase,
			updateTaskStatus: updateTaskStatus,
			deleteTask: deleteTask,
		}
}


type AssigneTaskRequest struct{
	Title string `json:"title"`
	Description string `json:"description"`
	AssignedTo int64 `json:"assigned_to"`
	Priority task.TaskPriority `json:"priority"`
	DueDate time.Time `json:"due_date"`
}


func (h *TaskHandler) AssignTaskToEmployee(c *fiber.Ctx) error{
	var req AssigneTaskRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request Body"})
	}
	userId, ok := c.Locals("user_id").(int64)
	if !ok{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user unauthorized"})
	}

	authenticatedUser, appErr := h.getUserByIdUseCase.Execute(c.Context(), userId)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}
	
	if authenticatedUser.Position == user.UserEmployee {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user Not Autherized to assign a task"})
	}

	// Get assigned user's name
	assignedToUser, appErr := h.getUserByIdUseCase.Execute(c.Context(), req.AssignedTo)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}
	assignedByUser := authenticatedUser.FirstName+" "+authenticatedUser.LastName 
	assignedToName := assignedToUser.FirstName + " " + assignedToUser.LastName

	 err := h.assignTask.Execute(c.Context(),req.Title, req.Description, assignedByUser, assignedToName, userId, req.AssignedTo, task.TaskStatusPending, req.Priority, req.DueDate)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Faild to Assign a task"})
	}
	return c.JSON(fiber.StatusCreated)
}


func (h *TaskHandler) GetAssignedTasksByEmployeeId(c *fiber.Ctx) error{
	userId, ok := c.Locals("user_id").(int64)
	if !ok{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User Unautherized"})
	}

	authenticatedUser, appErr := h.getUserByIdUseCase.Execute(c.Context(), userId)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	tasks, appErr := h.getAssignedTaskByEmployeeId.Execute(c.Context(), userId, authenticatedUser.Position == user.UserManegers)
	if appErr != nil{
		return c.Status((appErr.Code.HTTPStatus())).JSON(fiber.Map{"error": appErr.Message})
	}

	return c.JSON(tasks)
}

type UpdateTaskStatusRequest struct{
	Status task.TaskStatus `json:"status"`
}

func (h *TaskHandler) UpdateTaskStatus(c *fiber.Ctx) error{
	userId, ok := c.Locals("user_id").(int64)
	if !ok{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user unauthorized"})
	}

	var req UpdateTaskStatusRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	taskIdStr := c.Params("taskId")
	taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task id"})
	}

	appErr := h.updateTaskStatus.Execute(c.Context(), taskId, userId, req.Status)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	return c.JSON(fiber.Map{"message": "task status updated"})
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error{
	userId, ok := c.Locals("user_id").(int64)
	if !ok{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user unauthorized"})
	}

	authenticatedUser, appErr := h.getUserByIdUseCase.Execute(c.Context(), userId)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	if authenticatedUser.Position != user.UserManegers {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "only managers can delete tasks"})
	}

	taskIdStr := c.Params("taskId")
	taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task id"})
	}

	appErr = h.deleteTask.Execute(c.Context(), taskId, userId)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	return c.JSON(fiber.Map{"message": "task deleted successfully"})
}
