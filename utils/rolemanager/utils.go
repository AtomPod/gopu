package rolemanager

import (
	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/models"
)

//AddUserIDPrimaryAPI 添加用户ID的私有API指令
func AddUserIDPrimaryAPI(
	uid string, groupName string,
	mgr RoleManager, conf config.RBAC) error {
	roles := conf.Roles

	var groupRole *config.RBACGroup
	for _, role := range roles {
		if role.Name == groupName {
			groupRole = &role
			break
		}
	}

	if groupRole == nil {
		return ErrRoleNotExists
	}

	role := models.Role{
		Name: uid,
	}

	for _, api := range groupRole.IDAPIS {
		role.Permissions = append(role.Permissions, models.Permission{
			API:    api.Path,
			Method: api.Method,
		})
	}

	_, err := mgr.CreateRole(&role)
	if err != nil {
		return err
	}
	return nil
}
