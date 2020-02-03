package v1

import (
	"github.com/gin-gonic/gin"
	forms "github.com/ngs24313/gopu/api/forms/rbac"
	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/rolemanager"
)

//RBAC is rbac api
type RBAC struct {
	RoleMgr rolemanager.RoleManager
}

//Register register handles
func (r *RBAC) Register(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		v1.POST("/role", r.CreateRole)
		v1.DELETE("/role/:name", r.DeleteRole)

		v1.GET("/role/:name/user", r.GetUserForRoleByName)
		v1.GET("/role/:name", r.GetRoleByName)
		v1.GET("/role", r.GetRoleList)

		v1.POST("/role/:name/user/:id", r.AppendRoleForUser)
		v1.DELETE("/role/:name/user/:id", r.DeleteRoleForUser)
	}
}

//CreateRole handles POST /v1/role
func (r *RBAC) CreateRole(c *gin.Context) {
	form := &forms.RoleForm{}
	if err := c.ShouldBind(&form); err != nil {
		replyBadRequest(c, "Some fields is invalid", err)
		return
	}

	role := &models.Role{
		Name:        form.Role,
		Permissions: make([]models.Permission, len(form.Permissions)),
	}

	for i, permission := range form.Permissions {
		role.Permissions[i].API = permission.API
		role.Permissions[i].Method = permission.Method
	}

	ok, err := r.RoleMgr.CreateRole(role)
	if err != nil {
		replyInternalError(c, err)
		return
	}

	if !ok {
		replyBadRequest(c, "The role already has the permission", nil)
		return
	}

	replyOK(c, nil)
}

//DeleteRole handles DELETE /v1/role/:name
func (r *RBAC) DeleteRole(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		replyBadRequest(c, "The role name cannot be empty", nil)
		return
	}

	ok, err := r.RoleMgr.DeleteRole(name)
	if err != nil {
		replyInternalError(c, err)
		return
	}

	if !ok {
		replyBadRequest(c, "The role does not exist", nil)
		return
	}

	replyOK(c, nil)
}

//GetRoleByName handles GET /v1/role/:name
func (r *RBAC) GetRoleByName(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		replyBadRequest(c, "The role name cannot be empty", nil)
		return
	}
	role, err := r.RoleMgr.GetRoleByName(name)
	if err != nil {
		replyInternalError(c, err)
		return
	}
	replyOK(c, role)
}

//GetRoleList handles GET /v1/role
func (r *RBAC) GetRoleList(c *gin.Context) {
	form := &forms.RoleListForm{}
	if err := c.ShouldBind(form); err != nil {
		replyBadRequest(c, "Some fields is invalid", err)
		return
	}

	if form.Page <= 0 {
		form.Page = 1
	}

	if form.PageSize <= 0 || form.PageSize > 64 {
		form.PageSize = 16
	}

	reply, err := r.RoleMgr.ListRole(&rolemanager.ListRoleParams{
		Offset: (form.Page - 1) * form.PageSize,
		Count:  form.PageSize,
	})

	if err != nil {
		replyInternalError(c, err)
		return
	}
	replyOK(c, reply)
}

//GetUserForRoleByName handles GET /role/:name/user
func (r *RBAC) GetUserForRoleByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		replyBadRequest(c, "The role name cannot be empty", nil)
		return
	}

	user, err := r.RoleMgr.GetUserForRole(name)
	if len(user) == 0 || err != nil {
		replyNotFound(c, "No user belongs to this role.", err)
		return
	}
	replyOK(c, user)
}

//AppendRoleForUser handles POST  /role/:name/user/:id
func (r *RBAC) AppendRoleForUser(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		replyBadRequest(c, "The role name cannot be empty", nil)
		return
	}

	uid := c.Param("id")
	if uid == "" {
		replyBadRequest(c, "The user id cannot be empty", nil)
		return
	}

	_, err := r.RoleMgr.AddRoleForUser(uid, name)
	if err != nil && err != rolemanager.ErrUserHasRole {
		if err == rolemanager.ErrRoleNotExists {
			replyBadRequest(c, err.Error(), nil)
		} else {
			replyInternalError(c, err)
		}
		return
	}

	replyOK(c, nil)
}

//DeleteRoleForUser handles DELETE /role/:name/user/:id
func (r *RBAC) DeleteRoleForUser(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		replyBadRequest(c, "The role name cannot be empty", nil)
		return
	}

	uid := c.Param("id")
	if uid == "" {
		replyBadRequest(c, "The user id cannot be empty", nil)
		return
	}

	_, err := r.RoleMgr.DelRoleForUser(uid, name)
	if err != nil {
		if err == rolemanager.ErrUserNotHaveRole {
			replyBadRequest(c, err.Error(), nil)
		} else {
			replyInternalError(c, err)
		}
		return
	}
	replyOK(c, nil)
}
