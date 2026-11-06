package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// adminRolePredicate matches users that have an administrator role.
// It expects the user ID as the $1 query argument.
const adminRolePredicate = `EXISTS (
	SELECT 1
	FROM users_roles ur
	JOIN roles r ON r.id = ur.role_id
	WHERE ur.user_id = $1
	  AND LOWER(TRIM(r.name)) IN ('administrador', 'admin')
)`

// fileBaseAccessPredicate determines whether a file aliased as "f" can be
// accessed by the user identified by $1, ignoring the project deletion state.
const fileBaseAccessPredicate = `(
	` + adminRolePredicate + `
	OR f.created_by = $1
	OR EXISTS (
		SELECT 1
		FROM files_users fu
		WHERE fu.file_id = f.id
		  AND fu.user_id = $1
	)
	OR EXISTS (
		SELECT 1
		FROM files_groups fg
		JOIN users_groups ug ON ug.group_id = fg.group_id
		WHERE fg.file_id = f.id
		  AND ug.user_id = $1
	)
)`

// fileVisibilityPredicate determines whether a file aliased as "f" is visible
// to the user identified by $1. A file is visible when the user is an
// administrator, uploaded it, received direct access or received access
// through one of their groups, and its project has not been deleted.
const fileVisibilityPredicate = `(
	` + fileBaseAccessPredicate + `
	AND NOT EXISTS (
		SELECT 1
		FROM projects p
		WHERE p.id = f.project_id
		  AND p.is_deleted = 1
	)
)`

const fileEditPredicate = `(
	` + fileVisibilityPredicate + `
	AND (
		` + adminRolePredicate + `
		OR f.created_by = $1
		OR EXISTS (
			SELECT 1
			FROM files_users fu
			WHERE fu.file_id = f.id
			  AND fu.user_id = $1
			  AND fu.access_type::text = 'editar'
		)
		OR EXISTS (
			SELECT 1
			FROM files_groups fg
			JOIN users_groups ug ON ug.group_id = fg.group_id
			WHERE fg.file_id = f.id
			  AND ug.user_id = $1
			  AND fg.access_type::text = 'editar'
		)
	)
)`

const fileDownloadPredicate = `(
	` + fileVisibilityPredicate + `
	AND (
		` + adminRolePredicate + `
		OR f.created_by = $1
		OR EXISTS (
			SELECT 1
			FROM files_users fu
			WHERE fu.file_id = f.id
			  AND fu.user_id = $1
			  AND fu.access_type::text IN ('descargar', 'editar')
		)
		OR EXISTS (
			SELECT 1
			FROM files_groups fg
			JOIN users_groups ug ON ug.group_id = fg.group_id
			WHERE fg.file_id = f.id
			  AND ug.user_id = $1
			  AND fg.access_type::text IN ('descargar', 'editar')
		)
	)
)`

const fileDeletePredicate = `(
	` + fileVisibilityPredicate + `
	AND (
		` + adminRolePredicate + `
		OR f.created_by = $1
	)
)`

// projectVisibilityPredicate determines whether a project aliased as "p" is
// visible to the user identified by $1. A project is visible when it has not
// been deleted and the user is an administrator, created it or has access to
// at least one of its files.
const projectVisibilityPredicate = `(
	p.is_deleted = 0
	AND (
		` + adminRolePredicate + `
		OR p.created_by = $1
		OR EXISTS (
			SELECT 1
			FROM files f
			WHERE f.project_id = p.id
			  AND ` + fileVisibilityPredicate + `
		)
	)
)`

// clientVisibilityPredicate determines whether a client aliased as "c" is
// visible to the user identified by $1. A client is visible when the user is
// an administrator, created it or has access to at least one file belonging to
// one of its projects.
const clientVisibilityPredicate = `(
	` + adminRolePredicate + `
	OR c.created_by = $1
	OR EXISTS (
		SELECT 1
		FROM projects p
		JOIN files f ON f.project_id = p.id
		WHERE p.client_id = c.id
		  AND ` + fileVisibilityPredicate + `
	)
)`

// userCanManageRBACPredicate determines whether the user identified by $1 can
// manage global roles and groups. Administrators always can; any other user
// needs the global roles.gestionar permission through one of their roles.
const userCanManageRBACPredicate = `(
	` + adminRolePredicate + `
	OR EXISTS (
		SELECT 1
		FROM users_roles ur
		JOIN roles_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
		  AND p.name = 'roles.gestionar'
	)
)`

func userCanManageRBAC(ctx context.Context, db *pgxpool.Pool, userID int) (bool, error) {
	var canManage bool
	err := db.QueryRow(ctx, `SELECT `+userCanManageRBACPredicate, userID).Scan(&canManage)
	if err != nil {
		return false, err
	}
	return canManage, nil
}
