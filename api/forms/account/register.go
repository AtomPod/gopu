package account

//RegisterForm 注册表单
type RegisterForm struct {
	Username        string `json:"username" form:"username" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"`
	Email           string `json:"email" form:"email" binding:"required,email"`
	RegisterCode    string `json:"register_code" form:"register_code" binding:"required,len=6,numeric"`
}

//RegisterCodeForm 注册码表单
type RegisterCodeForm struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}
