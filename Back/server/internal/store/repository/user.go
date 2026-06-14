package repository

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]dtos.UserResponse, error) {
	q := `SELECT u.id, u.username, u.email, u.password, u.created_at, u.updated_at,
		COALESCE((SELECT COUNT(*) FROM files f WHERE f.created_by = u.id), 0) as uploaded_files
		FROM "user" u`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []dtos.UserResponse{}
	for rows.Next() {
		var u models.User
		var uploadedFiles int
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &uploadedFiles)
		if err != nil {
			return nil, err
		}
		users = append(users, toUserResponse(&u, uploadedFiles))
	}
	return users, nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, searchParam string) ([]dtos.UserResponse, error) {
	if searchParam == "" {
		return []dtos.UserResponse{}, nil
	}
	searchPattern := "%" + searchParam + "%"
	q := `SELECT u.id, u.username, u.email, u.password, u.created_at, u.updated_at,
		COALESCE((SELECT COUNT(*) FROM files f WHERE f.created_by = u.id), 0) as uploaded_files
		FROM "user" u
		WHERE u.username ILIKE $1 OR u.email ILIKE $1`
	rows, err := r.db.Query(ctx, q, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []dtos.UserResponse{}
	for rows.Next() {
		var u models.User
		var uploadedFiles int
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &uploadedFiles)
		if err != nil {
			return nil, err
		}
		users = append(users, toUserResponse(&u, uploadedFiles))
	}
	return users, nil
}

func (r *UserRepository) CreateUsers(user *models.User, ctx context.Context) (*models.User, error) {
	q := `INSERT INTO "user"(username, email, password) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := r.db.QueryRow(ctx, q, user.Username, user.Email, user.Password).Scan(&userId)
	if err != nil {
		return nil, err
	}
	user.ID = userId
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*dtos.UserResponse, error) {
	q := `SELECT u.id, u.username, u.email, u.password, u.created_at, u.updated_at,
		COALESCE((SELECT COUNT(*) FROM files f WHERE f.created_by = u.id), 0) as uploaded_files
		FROM "user" u WHERE u.id = $1`
	row := r.db.QueryRow(ctx, q, id)
	var u models.User
	var uploadedFiles int
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &uploadedFiles)
	if err != nil {
		return nil, err
	}
	resp := toUserResponse(&u, uploadedFiles)
	return &resp, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	q := `SELECT id, username, email, password, created_at, updated_at FROM "user" WHERE username = $1`
	row := r.db.QueryRow(ctx, q, username)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `SELECT id, username, email, password, created_at, updated_at FROM "user" WHERE email = $1`
	row := r.db.QueryRow(ctx, q, email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, req *dtos.UpdateUserRequest, id int) (*dtos.UserResponse, error) {
	q := `UPDATE "user" SET username = $1, email = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, q, req.Username, req.Email, time.Now(), id)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, id)
}
func (r *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	q := `DELETE FROM "user" WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func toUserResponse(u *models.User, uploadedFiles int) dtos.UserResponse {
	return dtos.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		// Password:      u.Password,
		CreatedBy:     u.Username,
		UpdatedBy:     u.Username,
		UploadedFiles: uploadedFiles,
		CreatedAt:     utils.ToCustomTime(u.CreatedAt),
		UpdatedAt:     utils.ToCustomTime(u.UpdatedAt),
	}
}

func (r *UserRepository) GetHashedPassword(ctx context.Context, userID int) (string, error) {
	q := `SELECT password FROM "user" WHERE id = $1`
	var hashedPassword string
	err := r.db.QueryRow(ctx, q, userID).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newHashedPassword string) error {
	q := `UPDATE "user" SET password = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, q, newHashedPassword, time.Now(), userID)
	return err
}
