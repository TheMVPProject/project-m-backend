package v1

import (
	"project_m_backend/app/user/usecase"
	"project_m_backend/pkg/auth/jwt"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/domain/user"
)


type UserHandler struct{
	createUserUserCase *usecase.CreatUserUseCase
	getUserByIdUserCase *usecase.GetUserByIdUseCase
	loginWithEmailUserCase *usecase.LoginWithEmailUseCase
	jwtService *jwt.JWTService
	accessTokenExpiry time.Duration
	refreshTokenExpiry time.Duration
}

func NewUserHandler(
	createUserUserCase *usecase.CreatUserUseCase,
	getUserByIdUseCase *usecase.GetUserByIdUseCase,
	loginWithEmailUseCase *usecase.LoginWithEmailUseCase,
	jwtService *jwt.JWTService,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
) *UserHandler{
	return  &UserHandler{
		createUserUserCase: createUserUserCase,
		getUserByIdUserCase: getUserByIdUseCase,
		loginWithEmailUserCase: loginWithEmailUseCase,
		jwtService: jwtService,
		accessTokenExpiry: accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}


func (h *UserHandler) GetUserById(c *fiber.Ctx) error{
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user Id"})
	}

	user, appErr := h.getUserByIdUserCase.Execute(c.Context(), id)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	return  c.JSON(user)
}


type CreateUserRequest struct{
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Position user.UserType `json:"position"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx)error{
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if appErr := h.createUserUserCase.Execute(c.Context(), req.FirstName, req.LastName, req.Email, req.Password, req.Position);
	appErr != nil{
		return  c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}
	
	loggedInUser, appErr := h.loginWithEmailUserCase.Execute(c.Context(), req.Email, req.Password)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	domainUser := appModel.ToDomainUser(loggedInUser)
	accessToken, err := h.jwtService.GenerateAccessToken(domainUser, h.accessTokenExpiry)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Faild to ganerate Access token"})
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(domainUser, h.refreshTokenExpiry)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Faild to Ganerate refresh toke"})
	}

	accessTokenExpiryTime := time.Now().Add(h.accessTokenExpiry)
	refreshTokenExpiryTime := time.Now().Add(h.refreshTokenExpiry)

	return c.JSON(
		LoginResponse{
			AccessToken: accessToken,
			AccessTokenExpiry: accessTokenExpiryTime,
			RefreshToken: refreshToken,
			RefreshTokenExpiry: refreshTokenExpiryTime,
			UserId: loggedInUser.ID,
			Position: loggedInUser.Position,
		},
	)
}

type LoginWithEmailRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct{
	AccessToken string `json:"access_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
	UserId             int64     `json:"user_id"`
	Position           user.UserType    `json:"position"`
}


func (h *UserHandler) LoginWithEmail(c *fiber.Ctx) error{
	var req LoginWithEmailRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body request"})
	}

	loggedInUser, appErr := h.loginWithEmailUserCase.Execute(c.Context(), req.Email, req.Password)
	if appErr != nil{
		return c.Status(appErr.Code.HTTPStatus()).JSON(fiber.Map{"error": appErr.Message})
	}

	domainUser := appModel.ToDomainUser(loggedInUser)

	accessToken, err := h.jwtService.GenerateAccessToken(domainUser, h.accessTokenExpiry)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Faild to generate access token"})
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(domainUser, h.refreshTokenExpiry)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Faild to generate refresh toke"})
	}

	accessTokenExpiryTime := time.Now().Add(h.accessTokenExpiry)
	refreshTokenExpiryTime := time.Now().Add(h.refreshTokenExpiry)

	return c.JSON(
		LoginResponse{
			AccessToken: accessToken,
			AccessTokenExpiry: accessTokenExpiryTime,
			RefreshToken: refreshToken,
			RefreshTokenExpiry: refreshTokenExpiryTime,
			UserId: loggedInUser.ID,
			Position: loggedInUser.Position,
		},
	)
}