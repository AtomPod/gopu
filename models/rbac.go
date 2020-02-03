package models

//Permission role permission
type Permission struct {
	Role   string `json:"role"`
	API    string `json:"api"`
	Method string `json:"method"`
}

//Role role model
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}
