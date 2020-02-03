package account

import (
	"github.com/ngs24313/gopu/models"
)

//ListUserForm for list user query
type ListUserForm struct {
	Query           string `json:"query" form:"query"`
	Page            int    `json:"page"  form:"page"`
	PageSize        int    `json:"page_size" form:"page_size"`
	CreateTimeStart int64  `json:"create_time_start" form:"create_time_start"`
	CreateTimeEnd   int64  `json:"create_time_end" form:"create_time_end"`

	OrderBy  string `json:"order_by" form:"orderby"`
	SortType string `json:"sort_type" form:"sort_type"`
}

//ListUserResultForm for list user query result
type ListUserResultForm struct {
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	PageCount int            `json:"page_count"`
	Users     []*models.User `json:"users"`
}
