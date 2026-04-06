package repository

import (
	"net/http"
	"server/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) GetAllUsers(ctx echo.Context) ([]*models.User, error) {
	q := `SELECT id, username, email, password, created_by, updated_by FROM users`
	rows, err := r.db.Query(ctx.Request().Context(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*models.User{}
	for rows.Next() {
		var u models.User
		// Escaneamos las columnas en el orden del SELECT
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		// Agregamos el usuario al slice del struct
		users = append(users, &u)
	}
	return users, nil
}
func (r *UserRepository) CreateUsers(user *models.User, ctx echo.Context) error {
	q := `INSERT INTO users (username, email, password,created_by,updated_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx.Request().Context(), q, user.Username, user.Email, user.Password, user.CreatedBy, user.UpdatedBy)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, "usuario creado con exito")
}

func (r *UserRepository) GetUserByEmail(ctx echo.Context, email string) (*models.User, error) {
	q := `SELECT id, username, email, password, created_by, updated_by FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx.Request().Context(), q, email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedBy, &u.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
