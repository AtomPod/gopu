package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/rolemanager"
	"github.com/rs/xid"
)

type casbinRoleManager struct {
	enforcer *casbin.SyncedEnforcer
}

//NewCasbinRoleManager 创建casbin角色管理器
func NewCasbinRoleManager(e *casbin.SyncedEnforcer) rolemanager.RoleManager {
	return &casbinRoleManager{
		enforcer: e,
	}
}

func (m *casbinRoleManager) DeleteRole(name string) (bool, error) {
	return m.enforcer.DeleteRole(name)
}

func (m *casbinRoleManager) CreateRole(role *models.Role) (bool, error) {
	if role == nil {
		return false, fmt.Errorf("The argument [role] does not be nil")
	}

	for _, p := range role.Permissions {
		_, err := m.enforcer.AddPermissionForUser(role.Name, p.API, p.Method)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (m *casbinRoleManager) GetRoleByName(name string) (*models.Role, error) {

	role := &models.Role{
		Name: name,
	}

	permissions := m.enforcer.GetPermissionsForUser(name)

	for _, permission := range permissions {
		role.Permissions = append(role.Permissions, models.Permission{
			API:    permission[1],
			Method: permission[2],
		})
	}

	return role, nil
}

func (m *casbinRoleManager) ListRole(params *rolemanager.ListRoleParams) (*rolemanager.ListRoleReply, error) {
	subjs := m.enforcer.GetAllSubjects()

	roleNames := make([]string, 0)
	for _, sub := range subjs {
		if _, err := xid.FromString(sub); err != nil {
			roleNames = append(roleNames, sub)
		}
	}

	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	count := params.Count
	if count <= 0 || count >= 64 {
		count = 16
	}

	startIndex := offset
	endIndex := startIndex + count

	if startIndex >= len(roleNames) {
		startIndex = len(roleNames)
	}

	if endIndex >= len(roleNames) {
		endIndex = len(roleNames)
	}

	reply := &rolemanager.ListRoleReply{}
	reply.TotalCount = len(roleNames)
	if startIndex < endIndex {
		roleNames = roleNames[startIndex:endIndex]
		reply.Roles = make([]*models.Role, 0, len(roleNames))

		for _, name := range roleNames {
			role, err := m.GetRoleByName(name)
			if err != nil {
				return nil, err
			}
			if len(role.Permissions) > 0 {
				reply.Roles = append(reply.Roles, role)
			}
		}
	}

	return reply, nil
}

func (m *casbinRoleManager) AddRoleForUser(user string, role string) (bool, error) {
	if !m.enforcer.HasPolicy(role) {
		return false, rolemanager.ErrRoleNotExists
	}
	ok, err := m.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, rolemanager.ErrUserHasRole
	}
	return true, nil
}

func (m *casbinRoleManager) DelRoleForUser(user string, role string) (bool, error) {
	ok, err := m.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, rolemanager.ErrUserNotHaveRole
	}
	return true, nil
}

func (m *casbinRoleManager) HasRoleForUser(user string, role string) (bool, error) {
	return m.enforcer.HasRoleForUser(user, role)
}

func (m *casbinRoleManager) GetRoleForUser(user string) ([]string, error) {
	return m.enforcer.GetRolesForUser(user)
}

func (m *casbinRoleManager) GetUserForRole(role string) ([]string, error) {
	return m.enforcer.GetUsersForRole(role)
}

func (m *casbinRoleManager) Validate(user string, permission *models.Permission) (bool, error) {
	return m.enforcer.Enforce(user, permission.API, permission.Method)
}
