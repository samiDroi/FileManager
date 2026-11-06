package repository

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) CreateUserWithRole(ctx context.Context, username, email, password string, roleID, createdBy int) (*dtos.CreatedUserResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var role dtos.AssignedRoleResponse
	if err := tx.QueryRow(ctx,
		`SELECT id, name, COALESCE(description, '')
		 FROM roles
		 WHERE id = $1
		 FOR SHARE`,
		roleID,
	).Scan(&role.ID, &role.Name, &role.Description); err != nil {
		return nil, err
	}

	var userID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO "user" (username, email, password, created_by, updated_by)
		 VALUES ($1, $2, $3, $4, $4)
		 RETURNING id`,
		username, email, password, createdBy,
	).Scan(&userID); err != nil {
		return nil, err
	}

	var creatorUsername string
	if err := tx.QueryRow(ctx, `SELECT username FROM "user" WHERE id = $1`, createdBy).Scan(&creatorUsername); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO users_roles (user_id, role_id, created_by, updated_by)
		 VALUES ($1, $2, $3, $3)`,
		userID, role.ID, createdBy,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dtos.CreatedUserResponse{
		ID:            userID,
		Username:      username,
		Email:         email,
		CreatedBy:     creatorUsername,
		UpdatedBy:     creatorUsername,
		UploadedFiles: 0,
		Roles:         []dtos.AssignedRoleResponse{role},
	}, nil
}

// AssignAdminRoleIfMissing bootstraps the system by making the oldest user an
// administrator when no administrator has been assigned yet. On a fresh
// database this is the first registered user.
func (r *UserRepository) AssignAdminRoleIfMissing(ctx context.Context) error {
	var roleID int
	err := r.db.QueryRow(ctx,
		`SELECT id FROM roles
		 WHERE LOWER(TRIM(name)) IN ('administrador', 'admin')
		 ORDER BY id
		 LIMIT 1`,
	).Scan(&roleID)
	if err != nil {
		return nil
	}

	var adminAssignments int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users_roles WHERE role_id = $1`,
		roleID,
	).Scan(&adminAssignments); err != nil {
		return err
	}
	if adminAssignments > 0 {
		return nil
	}

	var oldestUserID int
	if err := r.db.QueryRow(ctx,
		`SELECT id FROM "user" WHERE is_deleted = 0 ORDER BY id LIMIT 1`,
	).Scan(&oldestUserID); err != nil {
		return nil
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO users_roles (user_id, role_id, created_by, updated_by)
		 VALUES ($1, $2, $1, $1)
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		oldestUserID, roleID,
	)
	return err
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
	roles, err := r.getUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	resp.Roles = roles
	return &resp, nil
}

func (r *UserRepository) getUserRoles(ctx context.Context, userID int) ([]dtos.AssignedRoleResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.name, COALESCE(r.description, '')
		 FROM users_roles ur
		 JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = $1
		 ORDER BY r.name`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]dtos.AssignedRoleResponse, 0)
	for rows.Next() {
		var role dtos.AssignedRoleResponse
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
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

func (r *UserRepository) UpdateUser(ctx context.Context, req *dtos.UpdateUserRequest, id, updatedBy int) (*dtos.UserResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx,
		`UPDATE "user" SET username = $1, email = $2, updated_at = $3, updated_by = $4 WHERE id = $5`,
		req.Username, req.Email, time.Now(), updatedBy, id,
	)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	if req.RoleID != nil {
		var role dtos.AssignedRoleResponse
		if err := tx.QueryRow(ctx,
			`SELECT id, name, COALESCE(description, '')
			 FROM roles WHERE id = $1 FOR SHARE`,
			*req.RoleID,
		).Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM users_roles WHERE user_id = $1`, id); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO users_roles (user_id, role_id, created_by, updated_by)
			 VALUES ($1, $2, $3, $3)`,
			id, role.ID, updatedBy,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
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
