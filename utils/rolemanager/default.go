package rolemanager

import (
	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/models"
)

var (
	defaultRoleMgr RoleManager
)

//SetRoleManager 设置默认角色管理器
func SetRoleManager(mgr RoleManager) {
	defaultRoleMgr = mgr
}

//GetRoleManager 获取默认角色管理器
func GetRoleManager() RoleManager {
	return defaultRoleMgr
}

//ApplyConfigToRoleManager 应用默认配置信息到角色管理器
func ApplyConfigToRoleManager(mgr RoleManager, conf *config.Config) error {
	rbac := conf.RBAC
	for _, g := range rbac.Roles {
		role := &models.Role{
			Name: g.Name,
		}
		for _, api := range g.APIS {
			role.Permissions = append(
				role.Permissions,
				models.Permission{
					API:    api.Path,
					Method: api.Method,
				})
		}

		_, err := mgr.CreateRole(role)
		if err != nil {
			return err
		}
	}
	return nil
}
