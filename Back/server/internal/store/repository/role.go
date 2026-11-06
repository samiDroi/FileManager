package repository

import (
	"context"
	"fmt"
	"server/internal/domain/dtos"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CanManage(ctx context.Context, userID int) (bool, error) {
	return userCanManageRBAC(ctx, r.db, userID)
}

func (r *RoleRepository) IsAdmin(ctx context.Context, userID int) (bool, error) {
	var isAdmin bool
	err := r.db.QueryRow(ctx, `SELECT `+adminRolePredicate, userID).Scan(&isAdmin)
	return isAdmin, err
}

func (r *RoleRepository) HasPermission(ctx context.Context, userID int, permission string) (bool, error) {
	var allowed bool
	err := r.db.QueryRow(ctx,
		`SELECT `+adminRolePredicate+` OR EXISTS (
			SELECT 1
			FROM users_roles ur
			JOIN roles_permissions rp ON rp.role_id = ur.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE ur.user_id = $1 AND p.name = $2
		)`,
		userID, permission,
	).Scan(&allowed)
	return allowed, err
}

func (r *RoleRepository) GetPermissions(ctx context.Context) ([]dtos.PermissionResponse, error) {
	q := `SELECT id, name, COALESCE(description, '')
	      FROM permissions
	      ORDER BY name`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]dtos.PermissionResponse, 0)
	for rows.Next() {
		var permission dtos.PermissionResponse
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (r *RoleRepository) GetRoles(ctx context.Context) ([]dtos.RoleResponse, error) {
	q := `SELECT id, name, COALESCE(description, ''), created_at, updated_at
	      FROM roles
	      ORDER BY LOWER(name)`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	roles := make([]dtos.RoleResponse, 0)
	roleIndex := make(map[int]int)
	for rows.Next() {
		var role dtos.RoleResponse
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &createdAt, &updatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		role.Permissions = make([]dtos.PermissionResponse, 0)
		role.CreatedAt = utils.ToCustomTime(createdAt)
		role.UpdatedAt = utils.ToCustomTime(updatedAt)
		roleIndex[role.ID] = len(roles)
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	if len(roles) == 0 {
		return roles, nil
	}

	permissionsQ := `SELECT rp.role_id, p.id, p.name, COALESCE(p.description, '')
	                 FROM roles_permissions rp
	                 JOIN permissions p ON p.id = rp.permission_id
	                 ORDER BY p.name`
	permissionRows, err := r.db.Query(ctx, permissionsQ)
	if err != nil {
		return nil, err
	}
	defer permissionRows.Close()

	for permissionRows.Next() {
		var roleID int
		var permission dtos.PermissionResponse
		if err := permissionRows.Scan(&roleID, &permission.ID, &permission.Name, &permission.Description); err != nil {
			return nil, err
		}
		if index, ok := roleIndex[roleID]; ok {
			roles[index].Permissions = append(roles[index].Permissions, permission)
		}
	}
	return roles, permissionRows.Err()
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, roleID int) (*dtos.RoleResponse, error) {
	q := `SELECT id, name, COALESCE(description, ''), created_at, updated_at
	      FROM roles
	      WHERE id = $1`

	var role dtos.RoleResponse
	var createdAt time.Time
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, q, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	role.CreatedAt = utils.ToCustomTime(createdAt)
	role.UpdatedAt = utils.ToCustomTime(updatedAt)
	role.Permissions = make([]dtos.PermissionResponse, 0)

	permissionsQ := `SELECT p.id, p.name, COALESCE(p.description, '')
	                 FROM roles_permissions rp
	                 JOIN permissions p ON p.id = rp.permission_id
	                 WHERE rp.role_id = $1
	                 ORDER BY p.name`
	rows, err := r.db.Query(ctx, permissionsQ, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission dtos.PermissionResponse
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.Description); err != nil {
			return nil, err
		}
		role.Permissions = append(role.Permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) CreateRole(ctx context.Context, req dtos.RoleRequest, userID int) (*dtos.RoleResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var roleID int
	err = tx.QueryRow(ctx,
		`INSERT INTO roles (name, description, created_by, updated_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		req.Name, req.Description, userID, userID,
	).Scan(&roleID)
	if err != nil {
		return nil, err
	}

	if err := insertRolePermissions(ctx, tx, roleID, req.PermissionIDs, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetRoleByID(ctx, roleID)
}

func (r *RoleRepository) UpdateRole(ctx context.Context, roleID int, req dtos.RoleRequest, userID int) (*dtos.RoleResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx,
		`UPDATE roles
		 SET name = $1, description = $2, updated_by = $3, updated_at = now()
		 WHERE id = $4`,
		req.Name, req.Description, userID, roleID,
	)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("role not found")
	}

	if _, err := tx.Exec(ctx, `DELETE FROM roles_permissions WHERE role_id = $1`, roleID); err != nil {
		return nil, err
	}
	if err := insertRolePermissions(ctx, tx, roleID, req.PermissionIDs, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetRoleByID(ctx, roleID)
}

func (r *RoleRepository) DeleteRole(ctx context.Context, roleID int) error {
	result, err := r.db.Exec(ctx, `DELETE FROM roles WHERE id = $1`, roleID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

func (r *RoleRepository) GetUserPermissions(ctx context.Context, userID int) (*dtos.UserPermissionsResponse, error) {
	isAdmin, err := r.IsAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleRows, err := r.db.Query(ctx,
		`SELECT DISTINCT r.name
		 FROM users_roles ur
		 JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = $1
		 ORDER BY r.name`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	for roleRows.Next() {
		var role string
		if err := roleRows.Scan(&role); err != nil {
			roleRows.Close()
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := roleRows.Err(); err != nil {
		roleRows.Close()
		return nil, err
	}
	roleRows.Close()

	q := `SELECT DISTINCT p.name
	      FROM users_roles ur
	      JOIN roles_permissions rp ON rp.role_id = ur.role_id
	      JOIN permissions p ON p.id = rp.permission_id
	      WHERE ur.user_id = $1
	      ORDER BY p.name`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]string, 0)
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &dtos.UserPermissionsResponse{
		UserID:      userID,
		IsAdmin:     isAdmin,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func insertRolePermissions(ctx context.Context, tx pgx.Tx, roleID int, permissionIDs []int, userID int) error {
	seen := make(map[int]struct{}, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		if permissionID <= 0 {
			continue
		}
		if _, exists := seen[permissionID]; exists {
			continue
		}
		seen[permissionID] = struct{}{}

		_, err := tx.Exec(ctx,
			`INSERT INTO roles_permissions (role_id, permission_id, created_by, updated_by)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (role_id, permission_id) DO NOTHING`,
			roleID, permissionID, userID, userID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
