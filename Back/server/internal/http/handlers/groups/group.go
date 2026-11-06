package groups

import (
	"errors"
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	service     *services.GroupService
	fileService *services.FileService
	roleRepo    *repository.RoleRepository
}

func NewGroupHandler(service *services.GroupService, fileService *services.FileService, roleRepo *repository.RoleRepository) *GroupHandler {
	return &GroupHandler{
		service:     service,
		fileService: fileService,
		roleRepo:    roleRepo,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	groupRepo := repository.NewGroupRepository(db)
	fileRepo := repository.NewFileRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	groupSvc := services.NewGroupService(groupRepo)
	fileSvc := services.NewFileService(fileRepo)
	h := NewGroupHandler(groupSvc, fileSvc, roleRepo)

	e.GET("/groups", h.GetAll, middlewares.IsAuthenticated)
	e.POST("/groups/new", h.Create, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "roles", "gestionar"))
	e.PUT("/groups/edit/:id", h.Update, middlewares.IsAuthenticated)
	e.DELETE("/groups/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "roles", "gestionar"))
	e.GET("/groups/:id/users", h.GetUsers, middlewares.IsAuthenticated)
	e.GET("/groups/:id/available-users", h.GetAvailableUsers, middlewares.IsAuthenticated)
	e.POST("/groups/:id/users", h.AddUser, middlewares.IsAuthenticated)
	e.DELETE("/groups/:id/users/:userId", h.RemoveUser, middlewares.IsAuthenticated)
	e.GET("/groups/:id/files", h.GetFiles, middlewares.IsAuthenticated)
	e.POST("/groups/:id/files", h.AddFile, middlewares.IsAuthenticated)
	e.DELETE("/groups/:id/files/:fileId", h.RevokeFile, middlewares.IsAuthenticated)
}

func currentUserID(c echo.Context) (int, error) {
	userSession := middlewares.GetRequestUserSession(c)
	if !userSession.IsAuthenticated {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Usuario no autenticado")
	}
	return userSession.UserID, nil
}

func parseID(c echo.Context, param string) (int, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}
	return id, nil
}

func groupError(err error) error {
	switch {
	case errors.Is(err, services.ErrGroupAccessDenied):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrInvalidGroupData),
		errors.Is(err, services.ErrInvalidGroupRole):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *GroupHandler) GetAll(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groups, err := h.service.GetAllGroups(c.Request().Context(), userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, groups)
}

func (h *GroupHandler) Create(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	var req dtos.GroupRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	group, err := h.service.CreateGroup(c.Request().Context(), req, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusCreated, group)
}

func (h *GroupHandler) Update(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	var req dtos.GroupRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	group, err := h.service.UpdateGroup(c.Request().Context(), groupID, req, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) Delete(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.DeleteGroup(c.Request().Context(), groupID, userID); err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Grupo eliminado correctamente"})
}

func (h *GroupHandler) GetUsers(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	users, err := h.service.GetGroupUsers(c.Request().Context(), groupID, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *GroupHandler) GetAvailableUsers(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	users, err := h.service.GetAvailableUsers(c.Request().Context(), groupID, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *GroupHandler) AddUser(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	var req dtos.AddGroupUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	groupUser, err := h.service.AddGroupUser(c.Request().Context(), groupID, req, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusCreated, groupUser)
}

func (h *GroupHandler) RemoveUser(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}
	targetUserID, err := parseID(c, "userId")
	if err != nil {
		return err
	}

	if err := h.service.RemoveGroupUser(c.Request().Context(), groupID, targetUserID, userID); err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Usuario removido del grupo"})
}

func (h *GroupHandler) GetFiles(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	files, err := h.service.GetGroupFiles(c.Request().Context(), groupID, userID)
	if err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, files)
}

func (h *GroupHandler) AddFile(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}

	allowed, err := h.service.CanEditGroupFiles(c.Request().Context(), groupID, userID)
	if err != nil {
		return groupError(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden, services.ErrGroupAccessDenied.Error())
	}

	canGrant, err := h.roleRepo.HasPermission(c.Request().Context(), userID, "accesos.otorgar")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !canGrant {
		return echo.NewHTTPError(http.StatusForbidden, "No cuentas con el permiso accesos.otorgar")
	}

	var req dtos.AddGroupFileRequest
	if err := c.Bind(&req); err != nil || req.FileID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "file_id es obligatorio")
	}

	if err := h.fileService.GrantGroupAccess(c.Request().Context(), req.FileID, groupID, req.AccessType, userID); err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Archivo agregado al grupo correctamente"})
}

func (h *GroupHandler) RevokeFile(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	groupID, err := parseID(c, "id")
	if err != nil {
		return err
	}
	fileID, err := parseID(c, "fileId")
	if err != nil {
		return err
	}

	if err := h.service.RevokeGroupFile(c.Request().Context(), groupID, fileID, userID); err != nil {
		return groupError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Acceso al archivo revocado"})
}
