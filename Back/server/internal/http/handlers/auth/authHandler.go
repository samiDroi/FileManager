package auth

import (
	"fmt"
	"net/http"
	"server/internal/domain/models"
	"server/internal/http/api"
	"server/internal/http/middlewares"
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
	g.POST(paths.Auth.Logout.Route, h.Logout, middlewares.JWTMiddleware)
}

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
	user, err := h.service.RegisterUser(c, userDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al guardar en base de datos: " + err.Error(),
		})
	}

	// Generar JWT token to login
	token, err := createToken(c, user.ID, user.Username, user.Email)
	api.SetTokenCookie(c, token)
	//exit response with token
	return c.JSON(http.StatusCreated, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)

	//validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}
	fmt.Print(req.Username, req.Password)
	//call service to login the user
	user, err := h.service.LoginUser(c, req.Username, req.Password)
	if err != nil {
		fmt.Printf("DEBUG: Error en Scan de base de datos: %v\n", err) // <--- ESTO te dirá la verdad
		return err
	}

	token, err := createToken(c, user.ID, user.Username, user.Email)
	api.SetTokenCookie(c, token)
	//exit response with token
	return c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
func (h *AuthHandler) Logout(c echo.Context) error {
	api.ClearCookie(c)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sesión cerrada correctamente",
	})
}
func createToken(c echo.Context, user_ID int, username string, user_Email string) (string, error) {
	token, err := api.GenerateToken(user_ID, username, user_Email)
	if err != nil {
		return "", c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al generar token: " + err.Error(),
		})
	}
	return token, err
}
