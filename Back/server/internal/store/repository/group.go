package repository

import (
	"context"
	"fmt"
	"server/internal/domain/dtos"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) CanManage(ctx context.Context, userID int) (bool, error) {
	return userCanManageRBAC(ctx, r.db, userID)
}

func (r *GroupRepository) GetAllGroups(ctx context.Context, userID int, canManageAll bool) ([]dtos.GroupResponse, error) {
	q := `SELECT g.id, g.name, COALESCE(g.description, ''), g.created_at,
	             COALESCE((SELECT COUNT(*) FROM files_groups fg WHERE fg.group_id = g.id), 0) AS file_count,
	             COALESCE((SELECT COUNT(*) FROM users_groups ug WHERE ug.group_id = g.id), 0) AS user_count,
	             COALESCE((SELECT ug.role_in_group::text FROM users_groups ug WHERE ug.group_id = g.id AND ug.user_id = $1), '') AS role_in_group,
	             ($2 OR EXISTS (
	                 SELECT 1 FROM users_groups ug
	                 WHERE ug.group_id = g.id AND ug.user_id = $1 AND ug.role_in_group = 'admin'
	             )) AS can_manage
	      FROM public."groups" g
	      WHERE $2 OR EXISTS (
	          SELECT 1 FROM users_groups ug WHERE ug.group_id = g.id AND ug.user_id = $1
	      )
	      ORDER BY LOWER(g.name)`

	rows, err := r.db.Query(ctx, q, userID, canManageAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]dtos.GroupResponse, 0)
	for rows.Next() {
		var group dtos.GroupResponse
		var createdAt time.Time
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&createdAt,
			&group.FileCount,
			&group.UserCount,
			&group.RoleInGroup,
			&group.CanManage,
		); err != nil {
			return nil, err
		}
		group.CreatedAt = utils.ToCustomTime(createdAt)
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) GetGroupByID(ctx context.Context, groupID, userID int) (*dtos.GroupResponse, error) {
	q := `SELECT g.id, g.name, COALESCE(g.description, ''), g.created_at,
	             COALESCE((SELECT COUNT(*) FROM files_groups fg WHERE fg.group_id = g.id), 0) AS file_count,
	             COALESCE((SELECT COUNT(*) FROM users_groups ug WHERE ug.group_id = g.id), 0) AS user_count,
	             COALESCE((SELECT ug.role_in_group::text FROM users_groups ug WHERE ug.group_id = g.id AND ug.user_id = $1), '') AS role_in_group,
	             (` + userCanManageRBACPredicate + ` OR EXISTS (
	                 SELECT 1 FROM users_groups ug
	                 WHERE ug.group_id = g.id AND ug.user_id = $1 AND ug.role_in_group = 'admin'
	             )) AS can_manage
	      FROM public."groups" g
	      WHERE g.id = $2`

	var group dtos.GroupResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, q, userID, groupID).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&createdAt,
		&group.FileCount,
		&group.UserCount,
		&group.RoleInGroup,
		&group.CanManage,
	)
	if err != nil {
		return nil, err
	}
	group.CreatedAt = utils.ToCustomTime(createdAt)
	return &group, nil
}

func (r *GroupRepository) GetUserGroupRole(ctx context.Context, groupID, userID int) (string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE((
		     SELECT ug.role_in_group::text
		     FROM users_groups ug
		     WHERE ug.group_id = $1 AND ug.user_id = $2
		 ), '')`,
		groupID, userID,
	).Scan(&role)
	return role, err
}

func (r *GroupRepository) CanManageGroup(ctx context.Context, groupID, userID int) (bool, error) {
	var allowed bool
	err := r.db.QueryRow(ctx,
		`SELECT `+userCanManageRBACPredicate+` OR EXISTS (
		     SELECT 1 FROM users_groups ug
		     WHERE ug.group_id = $2 AND ug.user_id = $1 AND ug.role_in_group = 'admin'
		 )`,
		userID, groupID,
	).Scan(&allowed)
	return allowed, err
}

func (r *GroupRepository) CanEditGroupFiles(ctx context.Context, groupID, userID int) (bool, error) {
	var allowed bool
	err := r.db.QueryRow(ctx,
		`SELECT `+userCanManageRBACPredicate+` OR EXISTS (
		     SELECT 1 FROM users_groups ug
		     WHERE ug.group_id = $2 AND ug.user_id = $1 AND ug.role_in_group IN ('admin', 'editor')
		 )`,
		userID, groupID,
	).Scan(&allowed)
	return allowed, err
}

func (r *GroupRepository) GetAvailableUsers(ctx context.Context, groupID int) ([]dtos.AvailableUserResponse, error) {
	q := `SELECT u.id, u.username, u.email
	      FROM "user" u
	      WHERE u.is_deleted = 0
	        AND NOT EXISTS (
	            SELECT 1 FROM users_groups ug
	            WHERE ug.group_id = $1 AND ug.user_id = u.id
	        )
	      ORDER BY LOWER(u.username)`

	rows, err := r.db.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dtos.AvailableUserResponse, 0)
	for rows.Next() {
		var user dtos.AvailableUserResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *GroupRepository) CreateGroup(ctx context.Context, req dtos.GroupRequest, userID int) (*dtos.GroupResponse, error) {
	var groupID int
	err := r.db.QueryRow(ctx,
		`INSERT INTO public."groups" (name, description, created_by, updated_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		req.Name, req.Description, userID, userID,
	).Scan(&groupID)
	if err != nil {
		return nil, err
	}
	return r.GetGroupByID(ctx, groupID, userID)
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, groupID int, req dtos.GroupRequest, userID int) (*dtos.GroupResponse, error) {
	result, err := r.db.Exec(ctx,
		`UPDATE public."groups"
		 SET name = $1, description = $2, updated_by = $3, updated_at = now()
		 WHERE id = $4`,
		req.Name, req.Description, userID, groupID,
	)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("group not found")
	}
	return r.GetGroupByID(ctx, groupID, userID)
}

func (r *GroupRepository) DeleteGroup(ctx context.Context, groupID int) error {
	result, err := r.db.Exec(ctx, `DELETE FROM public."groups" WHERE id = $1`, groupID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("group not found")
	}
	return nil
}

func (r *GroupRepository) GetGroupUsers(ctx context.Context, groupID int) ([]dtos.GroupUserResponse, error) {
	q := `SELECT u.id, u.username, u.email, ug.role_in_group::text, ug.added_at
	      FROM users_groups ug
	      JOIN "user" u ON u.id = ug.user_id
	      WHERE ug.group_id = $1
	        AND u.is_deleted = 0
	      ORDER BY u.username`

	rows, err := r.db.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dtos.GroupUserResponse, 0)
	for rows.Next() {
		var user dtos.GroupUserResponse
		var addedAt time.Time
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.RoleInGroup,
			&addedAt,
		); err != nil {
			return nil, err
		}
		user.AddedAt = utils.ToCustomTime(addedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *GroupRepository) AddGroupUser(ctx context.Context, groupID int, userID int, roleInGroup string, createdBy int) (*dtos.GroupUserResponse, error) {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users_groups (user_id, group_id, role_in_group, created_by, updated_by)
		 VALUES ($1, $2, $3, $4, $4)
		 ON CONFLICT (user_id, group_id) DO UPDATE
		 SET role_in_group = EXCLUDED.role_in_group,
		     updated_by = EXCLUDED.updated_by,
		     updated_at = now()`,
		userID, groupID, roleInGroup, createdBy,
	)
	if err != nil {
		return nil, err
	}

	q := `SELECT u.id, u.username, u.email, ug.role_in_group::text, ug.added_at
	      FROM users_groups ug
	      JOIN "user" u ON u.id = ug.user_id
	      WHERE ug.group_id = $1
	        AND ug.user_id = $2`

	var user dtos.GroupUserResponse
	var addedAt time.Time
	err = r.db.QueryRow(ctx, q, groupID, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.RoleInGroup,
		&addedAt,
	)
	if err != nil {
		return nil, err
	}
	user.AddedAt = utils.ToCustomTime(addedAt)
	return &user, nil
}

func (r *GroupRepository) RemoveGroupUser(ctx context.Context, groupID int, userID int) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM users_groups WHERE group_id = $1 AND user_id = $2`,
		groupID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("group member not found")
	}
	return nil
}

func (r *GroupRepository) GetGroupFiles(ctx context.Context, groupID int) ([]dtos.GroupFileResponse, error) {
	q := `SELECT f.id, f.tittle, COALESCE(p.name, ''), COALESCE(ft.name, ''),
	             fg.access_type::text, fg.granted_at
	      FROM files_groups fg
	      JOIN files f ON f.id = fg.file_id
	      LEFT JOIN projects p ON p.id = f.project_id
	      LEFT JOIN file_types ft ON ft.id = f.type_id
	      WHERE fg.group_id = $1
	        AND f.is_deleted = 0
	        AND (p.id IS NULL OR p.is_deleted = 0)
	      ORDER BY LOWER(f.tittle)`

	rows, err := r.db.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]dtos.GroupFileResponse, 0)
	for rows.Next() {
		var file dtos.GroupFileResponse
		var grantedAt time.Time
		if err := rows.Scan(
			&file.ID,
			&file.Tittle,
			&file.ProjectName,
			&file.FileType,
			&file.AccessType,
			&grantedAt,
		); err != nil {
			return nil, err
		}
		file.GrantedAt = utils.ToCustomTime(grantedAt)
		files = append(files, file)
	}
	return files, rows.Err()
}

func (r *GroupRepository) AddGroupFile(ctx context.Context, groupID, fileID int, accessType string, grantedBy int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO files_groups (file_id, group_id, access_type, granted_by, created_by, updated_by)
		 VALUES ($1, $2, $3::file_access_enum, $4, $4, $4)
		 ON CONFLICT (file_id, group_id) DO UPDATE
		 SET access_type = EXCLUDED.access_type,
		     granted_by = EXCLUDED.granted_by,
		     granted_at = CURRENT_TIMESTAMP,
		     updated_at = CURRENT_TIMESTAMP,
		     updated_by = EXCLUDED.updated_by`,
		fileID, groupID, accessType, grantedBy,
	)
	return err
}

func (r *GroupRepository) RevokeGroupFile(ctx context.Context, groupID int, fileID int) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM files_groups WHERE group_id = $1 AND file_id = $2`,
		groupID, fileID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("file access not found")
	}
	return nil
}
