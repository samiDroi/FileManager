package repository

import (
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
	q := `SELECT id, username, email, password, created_by, updated_by FROM "user"`
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
func (r *UserRepository) CreateUsers(user *models.User, ctx echo.Context) (*models.User, error) {
	q := `INSERT INTO "user"(username, email, password) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := r.db.QueryRow(ctx.Request().Context(), q, user.Username, user.Email, user.Password).Scan(&userId)
	if err != nil {
		return nil, err
	}
	user.ID = userId
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx echo.Context, id int) (*models.User, error) {
	q := `SELECT id, username, email, password, created_by, updated_by FROM "user" WHERE id = $1`
	row := r.db.QueryRow(ctx.Request().Context(), q, id)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *UserRepository) GetUserByUsername(ctx echo.Context, username string) (*models.User, error) {
	q := `SELECT id, username, email, password, created_at, updated_at FROM "user" WHERE username = $1`
	row := r.db.QueryRow(ctx.Request().Context(), q, username)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByEmail(ctx echo.Context, email string) (*models.User, error) {
	q := `SELECT id, username, email, password, created_by, updated_by FROM "user" WHERE email = $1`
	row := r.db.QueryRow(ctx.Request().Context(), q, email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *UserRepository) UpdateUser(ctx echo.Context, user *models.User) error {
	q := `UPDATE "user" SET username = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx.Request().Context(), q, user.Username, user.Email, user.Password, user.UpdatedAt, user.ID)
	return err
}
