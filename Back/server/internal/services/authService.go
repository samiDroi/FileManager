package services

import (
	"errors"
	"fmt"
	"server/internal/store/repository"

	"server/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	ur *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{ur: repo}
}

func (s *AuthService) RegisterUser(ctx echo.Context, user models.User) (*models.User, error) {
	// Hash de la contraseña antes de guardar
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Le pedimos al store que lo guarde y retorne el usuario con ID asignado.
	createdUser, err := s.ur.CreateUsers(&user, ctx.Request().Context())
	if err != nil {
		return nil, err
	}

	if err := s.ur.AssignAdminRoleIfMissing(ctx.Request().Context()); err != nil {
		fmt.Printf("WARN: no se pudo asignar el rol Administrador al usuario inicial: %v\n", err)
	}

	return createdUser, nil
}

func (s *AuthService) LoginUser(ctx echo.Context, username string, password string) (*models.User, error) {
	// Buscamos el usuario por username
	user, err := s.ur.GetUserByUsername(ctx.Request().Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, echo.NewHTTPError(401, "Credenciales incorrectas")
		}
		return nil, echo.NewHTTPError(500, "Error interno del servidor")
	}

	// Validamos la contraseña con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, echo.NewHTTPError(401, "Contraseña incorrecta")
	}

	return user, nil
}
