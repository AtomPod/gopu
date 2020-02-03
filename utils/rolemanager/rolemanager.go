package rolemanager

import (
	"errors"

	"github.com/ngs24313/gopu/models"
)

var (
	//ErrRoleNotExists 角色不存在
	ErrRoleNotExists = errors.New("The role does not exists")
	//ErrUserHasRole 用户已经拥有该角色
	ErrUserHasRole = errors.New("The user already has the role ")
	//ErrUserNotHaveRole 用户没有拥有该角色权限
	ErrUserNotHaveRole = errors.New("The user does not have role")
)

//ListRoleParams 角色列表查询参数
type ListRoleParams struct {
	Offset int
	Count  int
}

//ListRoleReply 角色列表查询结果
type ListRoleReply struct {
	TotalCount int            `json:"total_count"`
	Roles      []*models.Role `json:"roles"`
}

//RoleManager role manager interface
type RoleManager interface {
	DeleteRole(name string) (bool, error)
	CreateRole(role *models.Role) (bool, error)
	GetRoleByName(name string) (*models.Role, error)
	ListRole(params *ListRoleParams) (*ListRoleReply, error)

	AddRoleForUser(user string, role string) (bool, error)
	DelRoleForUser(user string, role string) (bool, error)
	HasRoleForUser(user string, role string) (bool, error)
	GetRoleForUser(user string) ([]string, error)
	GetUserForRole(role string) ([]string, error)

	Validate(user string, permission *models.Permission) (bool, error)
}
