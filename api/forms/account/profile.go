package account

import "mime/multipart"

//ProfileForm  profile form
type ProfileForm struct {
	Avatar   *multipart.FileHeader `form:"avatar" binding:"omitempty"`
	Nickname string                `form:"nickname" json:"nickname" binding:"omitempty"`
	Company  string                `form:"company" json:"company" binding:"omitempty"`
	Location string                `form:"location" json:"location" binding:"omitempty"`
}
