package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	v1 "assignly/api/handler/v1"
	"assignly/api/middleware"

	userPorts "assignly/app/user/ports"
	taskPorts "assignly/app/task/ports"
	postgres "assignly/app/user/repository/postgres"
	taskPostgres "assignly/app/task/repository/postgres"
	"assignly/app/user/usecase"
	taskUsecase "assignly/app/task/usecase"
	"assignly/domain/user"
	"assignly/pkg/auth/jwt"
	"assignly/pkg/crypto"
	"assignly/pkg/security"
	"assignly/pkg/validator"
	"time"

	"github.com/gofiber/fiber/v2"
)



func main(){
	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	var userRepo userPorts.UserRepository
	userRepo, err = postgres.NewPostgresUserRepository(db)
	if err != nil{
		log.Fatal(err)
	}

	var taskRepo taskPorts.TaskRepository
	taskRepo, err = taskPostgres.NewPostgresTaskRepository(db)
	if err != nil{
		log.Fatal(err)
	}

	hasher := security.NewBcryptHasher()
	emailValidator := validator.NewEmailValidator()
	userFactory := user.NewUserFactory(hasher, emailValidator)
	// tokenGenerator := security.NewTokenGenerator()

	encryptor, err := crypto.NewAESEncryptor([]byte("12345678901234567890123456789012"))
	if err != nil{
		log.Fatalf("Faild to create encrypeter: %v", err)
	}

	jwtService := jwt.NewJWTService(os.Getenv("JWT_SECRET"), encryptor)

	createUserUseCase := usecase.NewCreateUserUseCase(userRepo, userFactory)
	getUserByIdUseCase := usecase.NewGetUserByIDUseCase(userRepo)
	loginWithEmail := usecase.NewLoginWithEmailUseCase(userRepo, hasher, userFactory)
	getEmployeeUseCase := usecase.NewGetEmployeeListUseCase(userRepo)

	assignTaskUseCase := taskUsecase.NewAssignTaskUseCase(taskRepo)
	getAssignedTaskByEmployeeIdUseCase := taskUsecase.NewGetAssignedTaskByEmployeeIdUseCase(taskRepo)
	updateTaskStatusUseCase := taskUsecase.NewUpdateTaskStatusUseCase(taskRepo)
	deleteTaskUseCase := taskUsecase.NewDeleteTaskUseCase(taskRepo)

	authMiddleWare := middleware.NewAuthMiddleware(jwtService)

	userHandler := v1.NewUserHandler(
		createUserUseCase,
		getUserByIdUseCase,
		loginWithEmail,
		getEmployeeUseCase,
		jwtService,
		time.Hour*24*7,
		time.Hour*24*365,
	)

	taskHandler := v1.NewTaskHandler(
		assignTaskUseCase,
		getAssignedTaskByEmployeeIdUseCase,
		getUserByIdUseCase,
		updateTaskStatusUseCase,
		deleteTaskUseCase,
	)

	app := fiber.New()

	users := app.Group("/users")
	users.Post("/create", userHandler.CreateUser)
	users.Post("/login", userHandler.LoginWithEmail)
	users.Post("/refresh-token", userHandler.RefreshToken)
	users.Get("/get", userHandler.GetUserById)
	users.Get("/employees-list", authMiddleWare, userHandler.GetEmployeeList)

	tasks := app.Group("/tasks")
	tasks.Post("/assign-task",authMiddleWare, taskHandler.AssignTaskToEmployee)
	tasks.Get("/get-assigned-task", authMiddleWare, taskHandler.GetAssignedTasksByEmployeeId)
	tasks.Put("/status/:taskId", authMiddleWare, taskHandler.UpdateTaskStatus)
	tasks.Delete("/delete/:taskId", authMiddleWare, taskHandler.DeleteTask)

	fmt.Println("Server listening on port 8080")
	log.Fatal(app.Listen(":8080"))
}