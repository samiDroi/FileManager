package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"server/internal/config"
	"server/internal/domain/dtos"
	"server/internal/http/api"
	"server/internal/store/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type resetPasswordHandler struct {
	repo *repository.UserRepository
	db   *pgxpool.Pool
}

func NewResetPasswordHandler(db *pgxpool.Pool) *resetPasswordHandler {
	return &resetPasswordHandler{
		repo: repository.NewUserRepository(db),
		db:   db,
	}
}

func (h *resetPasswordHandler) ForgotPassword(c echo.Context) error {
	input := new(dtos.UserForgotPasswordRequest)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, err := h.repo.GetUserByEmail(c.Request().Context(), input.Email)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Si el correo existe, se enviará un enlace de recuperación"})
	}

	token, err := api.GenerateResetPasswordToken(user.ID, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating reset token"})
	}

	if err := sendResetPasswordEmail(user.Email, token); err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error sending reset email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Si el correo existe, se enviará un enlace de recuperación"})
}

func (h *resetPasswordHandler) RecoveryPassword(c echo.Context) error {
	input := new(dtos.UserForgotPassword)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	claims := &api.Claims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado")
		}

		ctx := c.Request().Context()
		currentHash, err := h.repo.GetHashedPassword(ctx, claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("usuario no válido")
		}

		return []byte(config.Config.SecretKey + currentHash), nil
	})

	if err != nil || !token.Valid {
		log.Printf("JWT validation failed for token: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "El enlace expiró o ya fue utilizado"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password for user %d: %v", claims.UserID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error interno"})
	}

	ctx := c.Request().Context()
	if err := h.repo.UpdatePassword(ctx, claims.UserID, string(hashedPassword)); err != nil {
		log.Printf("Error updating password for user %d: %v", claims.UserID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No se pudo actualizar la contraseña"})
	}

	api.ClearCookies(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Contraseña reestablecida correctamente"})
}

func sendResetPasswordEmail(email, token string) error {
	cfg := config.Config
	resetLink := fmt.Sprintf("%s/recovery-password?token=%s", cfg.FrontendURL, token)

	subject := "Restablece tu contraseña - CattFileManager"
	body := fmt.Sprintf(`Hola,

Recibimos una solicitud para restablecer la contraseña de tu cuenta.

Haz clic en el siguiente enlace para restablecerla (válido por 15 minutos):

%s

Si no solicitaste este cambio, ignora este correo.

Saludos,
El equipo de CattFileManager`, resetLink)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.MailFromAddress, email, subject, body)

	addr := fmt.Sprintf("%s:%s", cfg.MailHost, cfg.MailPort)
	auth := smtp.PlainAuth("", cfg.MailUsername, cfg.MailPassword, cfg.MailHost)

	return smtp.SendMail(addr, auth, cfg.MailFromAddress, []string{email}, []byte(msg))
}
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	input := new(dtos.UserLoggedResetPwd)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	// userID := api.GetUserIDFromContext(c)

	// if err := h.service.ChangePassword(c.Request().Context(), userID, input.CurrentPassword, input.NewPassword); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	// }

	api.ClearCookies(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Contraseña cambiada correctamente. Por favor, inicia sesión de nuevo."})
}
