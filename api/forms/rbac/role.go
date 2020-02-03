package rbac

//PermissionForm permission http form
type PermissionForm struct {
	API    string `json:"api" form:"api" binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
}

//RoleForm role http form
type RoleForm struct {
	Role        string           `json:"role" form:"role" binding:"required,alphanum,ge=1,lt=20"`
	Permissions []PermissionForm `json:"permissions" form:"permissions" binding:"required"`
}

//RoleListForm role list http form
type RoleListForm struct {
	Page     int `json:"page" form:"page" binding:"omitempty"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty"`
}

//RoleUpdateForm role update htto form
type RoleUpdateForm struct {
	API    string `json:"api" form:"api" binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
}
