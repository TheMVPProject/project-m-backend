package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	userPorts "project_m_backend/app/user/ports"
	postgres "project_m_backend/app/user/repository/postgres"
	"project_m_backend/app/user/usecase"
	"project_m_backend/domain/user"
	"project_m_backend/pkg/auth/jwt"
	"project_m_backend/pkg/crypto"
	"project_m_backend/pkg/security"
	"project_m_backend/pkg/validator"
	"time"
	"github.com/gofiber/fiber/v2"
	v1 "project_m_backend/api/handler/v1"
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

	hasher := security.NewBcryptHasher()
	emailValidator := validator.NewEmailValidator()
	userFactory := user.NewUserFactory(hasher, emailValidator)
	// tokenGenerator := security.NewTokenGenerator()

	encryptor, err := crypto.NewAESEncryptor([]byte("12345678901234567890123456789012"))
	if err != nil{
		log.Fatalf("Faild to create encrypeter: %v", err)
	}

	jwtService := jwt.NewJWTService("your-secret-key", encryptor)

	createUserUseCase := usecase.NewCreateUserUseCase(userRepo, userFactory)
	getUserByIdUseCase := usecase.NewGetUserByIDUseCase(userRepo)
	loginWithEmail := usecase.NewLoginWithEmailUseCase(userRepo, hasher, userFactory)


	userHandler := v1.NewUserHandler(
		createUserUseCase,
		getUserByIdUseCase,
		loginWithEmail,
		jwtService,
		time.Minute*15,
		time.Hour*24*7,
	)

	app := fiber.New()

	users := app.Group("/users")
	users.Post("/create", userHandler.CreateUser)
	users.Post("/login", userHandler.LoginWithEmail)
	users.Get("/get", userHandler.GetUserById)

	fmt.Println("Server listening on port 8080")
	log.Fatal(app.Listen(":8080"))
}