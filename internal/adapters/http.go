package adapters

import (
	"log/slog"
	"net/http"

	"github.com/dkar-dev/hitpipe/internal/service"
	"github.com/labstack/echo/v4"
)

type UserAPI struct {
	userService *service.UserService
	log         *slog.Logger
}

func NewUserAPI(s *service.UserService, log *slog.Logger) *UserAPI {
	return &UserAPI{
		userService: s,
		log:         log,
	}
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Email    string `json:"email" validate:"required,email"`
//Password string `json:"password" validate:"required,min=8"`

func (a *UserAPI) Register(c echo.Context) error {
	var req RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	//if err := c.Validate(req); err != nil {
	//	return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	//}

	aggregate, err := a.userService.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"success": "true", "id": aggregate.User.ID})
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *UserAPI) Login(c echo.Context) error {
	var req LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	aggregate, err := a.userService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"success": "true", "id": aggregate.User.ID})
}

func (a *UserAPI) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")

	//log.Info("token verification start", token)

	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "token is empty"})
	}

	err := a.userService.VerifyEmail(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "verified"})
}
