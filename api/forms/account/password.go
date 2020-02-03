package account

//ResetPasswordForm for reset password
type ResetPasswordForm struct {
	OldPassword     string `json:"old_password" form:"old_password" binding:"omitempty"`
	NewPassword     string `json:"new_password" form:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=NewPassword"`
	ResetCode       string `json:"reset_code" form:"reset_code" binding:"omitempty,len=6,numeric"`
}

//PasswordResetCodeForm for create password reset code
type PasswordResetCodeForm struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}
