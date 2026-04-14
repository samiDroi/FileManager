package services

import (
	"database/sql"
	"fmt"
	"server/internal/store/repository"

	"server/internal/domain/models"

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
	return s.ur.CreateUsers(&user, ctx)
}

func (s *AuthService) LoginUser(ctx echo.Context, username string, password string) (*models.User, error) {
	// Buscamos el usuario por username
	// fmt.Print(ctx, username)
	user, err := s.ur.GetUserByUsername(ctx, username)
	if err != nil {
		// Imprime el error exacto en tu consola para que puedas verlo
		fmt.Printf("ERROR: %v\n", err)

		// Si quieres ser más preciso con la respuesta HTTP:
		if err == sql.ErrNoRows {
			return nil, echo.NewHTTPError(401, "Credenciales incorrectas")
		}
		return nil, echo.NewHTTPError(500, "Error interno del servidor")
	}
	// if err != nil {
	// 	return nil, echo.NewHTTPError(401, fmt.Sprintf("usuario %s no encontrado", username))
	// }

	// Validamos la contraseña con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, echo.NewHTTPError(401, "Contraseña incorrecta")
	}

	return user, nil
}
