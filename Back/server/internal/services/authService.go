package services

import (
	"server/internal/store/repository"

	"server/internal/domain/models"

	"github.com/labstack/echo/v4"
)

type AuthService struct {
	ur *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{ur: repo}
}
func (s *AuthService) RegisterUser(ctx echo.Context, user models.User) error {
	// Aquí es donde vive la lógica.
	// Por ahora, solo le pedimos al store que lo guarde.
	return s.ur.CreateUsers(&user, ctx)
}

func (s *AuthService) LoginUser(ctx echo.Context, email string, password string) (*models.User, error) {
	// Buscamos el usuario por email
	user, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Validamos la contraseña (comparación simple - en producción usar bcrypt)
	if user.Password != password {
		return nil, echo.NewHTTPError(401, "Contraseña incorrecta")
	}

	return user, nil
}
