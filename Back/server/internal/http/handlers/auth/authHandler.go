package auth

import (
	"net/http"
	// "server/internal/http/handlers/auth/dto"
	"server/internal/domain/models"
	"server/internal/paths"
	"server/internal/services"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}
func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	// A. Construimos la cadena de mando de este módulo
	repo := repository.NewUserRepository(db)
	svc := services.NewAuthService(repo)
	h := &AuthHandler{service: svc}

	// B. Agrupamos y definimos las rutas como a ti te gusta
	g := e.Group(paths.Auth.Prefix) // o usar tu paths.Auth.Prefix

	// Si tuvieras middlewares, van aquí: g.Use(middleware.Logger())

	g.POST(paths.Auth.Register.Route, h.Register)
	g.POST(paths.Auth.Login.Route, h.Login)
}

// func (h *AuthHandler) Login(c echo.Context) error {

// }
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)

	//validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}
	//Adapt the request dto to the user model
	userDto := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	//call service to register the user
	err := h.service.RegisterUser(c, userDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al guardar en base de datos: " + err.Error(),
		})
	}

	//exit response
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "¡Usuario registrado correctamente!",
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(RegisterRequest)

	//validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}

	//call service to login the user
	user, err := h.service.LoginUser(c, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Credenciales inválidas: " + err.Error(),
		})
	}

	//exit response with user data
	return c.JSON(http.StatusOK, user)
}
