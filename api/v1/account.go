package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/middleware"
	"github.com/ngs24313/gopu/utils"
	"github.com/ngs24313/gopu/utils/log"
	"github.com/ngs24313/gopu/utils/mailer"
	"github.com/ngs24313/gopu/utils/rolemanager"
	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
	db "github.com/ngs24313/gopu/api/database"
	apierr "github.com/ngs24313/gopu/api/error"
	forms "github.com/ngs24313/gopu/api/forms/account"
	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/cache"
	"github.com/ngs24313/gopu/utils/password"
	"go.uber.org/zap"
)

//Account is account api
type Account struct {
	ADB            db.AccountDatabase
	RoleMgr        rolemanager.RoleManager
	AuthMiddleware *middleware.Auth
	Config         config.Config
	Cache          cache.Cache
}

//Register register handles
func (a *Account) Register(router *gin.RouterGroup) {
	authMiddleware, err := a.AuthMiddleware.Middleware()
	if err != nil {
		panic(err)
	}

	v1 := router.Group("/v1")

	session := v1.Group("/session")
	{
		session.POST("/", authMiddleware.LoginHandler)
		session.DELETE("/", authMiddleware.LogoutHandler)
		session.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	{
		v1.POST("/user/password/reset_code", a.CreatePasswordResetCode)

		v1.POST("/user/register_code", a.CreateUserRegisterCode)

		v1.POST("/user", a.RegisterUser)

		v1.GET("/user/:id", a.GetUserByID)

		v1.PUT("/user/:id/password", a.ResetUserPassword)

		user := v1.Group("/")
		user.Use(authMiddleware.MiddlewareFunc())
		{
			user.GET("/current_user", a.CurrentUser)

			user.PUT("/user/:id/profile", a.UpdateUserProfile)

			user.DELETE("/user/:id", a.DeleteUser)

			user.GET("/user", a.ListUser)
		}
	}
}

//CreateUserRegisterCode handles POST /v1/user/register_code
func (a *Account) CreateUserRegisterCode(c *gin.Context) {
	form := &forms.RegisterCodeForm{}
	if err := c.ShouldBind(form); err != nil {
		replyBadRequest(c, "Some fields is not valid", err)
		return
	}

	if _, err := a.ADB.GetUserByEmail(c.Request.Context(), form.Email); err != nil {
		if err != db.ErrNotFound {
			replyInternalError(c, err)
			return
		}
	} else {
		replyBadRequest(c, db.ErrEmailAlreadyExists.Error(), nil)
		return
	}

	resetCode := utils.RandomDigit(6)
	codeKey := a.Config.Services.Account.RegisterCodePrefix + form.Email + "." + resetCode

	if err := a.Cache.Set(&cache.Entity{
		Key:        codeKey,
		Value:      []byte(""),
		Expiration: a.Config.Services.Account.RegisterCodeExpiration,
	}); err != nil {
		replyInternalError(c, err)
		return
	}

	replyEmail(c, mailer.User{
		Address: a.Config.Mailer.Username,
	}, mailer.User{
		Address: form.Email,
	}, a.Config.Mailer.EmailTemplates.RegisterCodeName,
		map[string]string{
			"code": resetCode,
		},
		nil,
	)
}

//RegisterUser handles POST /v1/user
func (a *Account) RegisterUser(c *gin.Context) {
	form := &forms.RegisterForm{}
	if err := c.ShouldBind(form); err != nil {
		replyBadRequest(c, "Some fields is not valid", err)
		return
	}

	if form.ConfirmPassword != form.Password {
		replyBadRequest(c, "Password and confirm password does not match", nil)
		return
	}

	codeKey := a.Config.Services.Account.RegisterCodePrefix + form.Email + "." + form.RegisterCode
	if _, err := a.Cache.Get(codeKey); err != nil {
		if err == cache.ErrNotFound {
			replyBadRequest(c, "Register code is invalid", nil)
			return
		}
		replyInternalError(c, err)
		return
	}

	user := &models.User{
		Username: form.Username,
		Password: password.GenHashPassword(form.Password),
		Email:    form.Email,
		Profile: models.Profile{
			Nickname: form.Username,
		},
	}

	if ok, err := a.ADB.UserIsExists(c.Request.Context(), user); ok || err != nil {
		if !ok && err != nil {
			replyInternalError(c, err)
			return
		} else if ok {
			replyBadRequest(c, err.Error(), nil)
			return
		}
	}

	createdUser, err := a.ADB.CreateUser(c.Request.Context(), user)
	if err != nil {
		replyInternalError(c, err)
		return
	}

	count, err := a.ADB.CountUser(c.Request.Context())
	if err != nil {
		replyInternalError(c, err)
		return
	}

	if count == 1 {
		if _, err := a.RoleMgr.AddRoleForUser(createdUser.ID, a.Config.RBAC.AdminName); err != nil {
			if err != rolemanager.ErrUserHasRole {
				log.Logger(c.Request.Context()).Warn("Failed to set user role", zap.Error(err))
			}
		}
	} else {
		if _, err := a.RoleMgr.AddRoleForUser(createdUser.ID, a.Config.RBAC.UserName); err != nil {
			if err != rolemanager.ErrUserHasRole {
				log.Logger(c.Request.Context()).Warn("Failed to set user role", zap.Error(err))
			}
		}

		if err := rolemanager.AddUserIDPrimaryAPI(createdUser.ID, a.Config.RBAC.UserName, a.RoleMgr, a.Config.RBAC); err != nil {
			log.Logger(c.Request.Context()).Warn("Failed to set user primary api", zap.Error(err))
		}
	}

	if err := a.userWithRoles(createdUser); err != nil {
		replyInternalError(c, err)
		return
	}

	if err := a.Cache.Del(codeKey); err != nil {
		log.Logger(c.Request.Context()).Warn(
			"Failed to delete register code",
			zap.String("key", codeKey),
			zap.Error(err),
		)
	}
	replyOK(c, createdUser)
}

//CurrentUser handles GET  /v1/current_user
func (a *Account) CurrentUser(c *gin.Context) {
	a.withUserByContext(c, func(user *models.User) {
		path := fmt.Sprintf("/v1/user/%s", user.ID)
		c.Redirect(http.StatusTemporaryRedirect, path)
	})
}

//DeleteUser handles DELETE /v1/user/:id
func (a *Account) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := a.ADB.DeleteUser(c.Request.Context(), id); err != nil {
		if err == db.ErrNotFound {
			replyError(c, apierr.NewAppError(
				http.StatusNotFound,
				fmt.Sprintf("user id [%s] is not found", id),
			))
			return
		}
		replyInternalError(c, err)
		return
	}

	replyOK(c, gin.H{
		"id": id,
	})
}

//GetUserByID handles GET /user/:id
func (a *Account) GetUserByID(c *gin.Context) {
	a.withUserByID(c, func(user *models.User) {
		if err := a.userWithRoles(user); err != nil {
			replyInternalError(c, err)
			return
		}
		replyOK(c, user)
	})
}

//ListUser handles GET /user
func (a *Account) ListUser(c *gin.Context) {

	form := forms.ListUserForm{}

	if err := c.ShouldBind(&form); err != nil {
		replyBadRequest(c, "Some fields is not valid", err)
		return
	}

	page := form.Page
	if page <= 0 {
		page = 1
	}
	pageSize := form.PageSize
	if pageSize <= 0 || pageSize > 32 {
		pageSize = 16
	}

	query := db.UserListQuery{
		Query:           form.Query,
		Offset:          (page - 1) * pageSize,
		Count:           pageSize,
		CreateTimeStart: time.Unix(form.CreateTimeStart, 0),
		CreateTimeEnd:   time.Unix(form.CreateTimeEnd, 0),
		WithCount:       true,
		OrderBy:         strings.Split(form.OrderBy, ","),
		SortType:        strings.Split(form.SortType, ","),
	}

	result, err := a.ADB.ListUser(c.Request.Context(), query)
	if err != nil {
		if err == db.ErrNotFound {
			replyError(c, apierr.NewAppError(http.StatusNotFound))
			return
		}
		replyInternalError(c, err)
		return
	}

	resultForm := forms.ListUserResultForm{
		Page:      page,
		PageSize:  pageSize,
		PageCount: int(result.Count) / pageSize,
		Users:     result.Users,
	}

	if int(result.Count)%pageSize != 0 {
		resultForm.PageCount++
	}

	for _, user := range resultForm.Users {
		if user != nil {
			if err := a.userWithRoles(user); err != nil {
				replyInternalError(c, err)
				return
			}
		}
	}
	replyOK(c, &resultForm)
}

//UpdateUserProfile handles PUT /user/:id/profile
func (a *Account) UpdateUserProfile(c *gin.Context) {
	form := &forms.ProfileForm{}
	if err := c.ShouldBind(form); err != nil {
		replyBadRequest(c, "Some fields is not valid", err)
		return
	}

	a.withUserByID(c, func(user *models.User) {
		avatar := user.Profile.Avatar
		if form.Avatar != nil {
			avatarID := xid.New().String()
			cfg := a.Config
			ext := filepath.Ext(form.Avatar.Filename)
			avatarID = avatarID + ext
			path := filepath.Join(cfg.AvatarBasePath(), avatarID)

			log.Logger(context.Background()).Debug("Upload file", zap.String("filepath", path))
			if err := c.SaveUploadedFile(form.Avatar, path); err != nil {
				replyInternalError(c, err)
				return
			}
			if avatar != "" {
				oldPath := filepath.Join(cfg.AvatarBasePath(), avatarID)
				if err := os.Remove(oldPath); err != nil {
					log.Logger(c.Request.Context()).Warn("Failed to remove avatar file", zap.String("path", oldPath), zap.Error(err))
				}
			}
			avatar = avatarID
		}

		if err := a.ADB.UpdateProfile(c.Request.Context(), &models.Profile{
			ID:       user.ProfileID,
			Avatar:   avatar,
			Company:  form.Company,
			Nickname: form.Nickname,
			Location: form.Location,
		}); err != nil {
			replyInternalError(c, err)
			return
		}

		replyOK(c, nil)
	})
}

//ResetUserPassword handles PUT /user/:id/password
func (a *Account) ResetUserPassword(c *gin.Context) {
	form := &forms.ResetPasswordForm{}
	if err := c.ShouldBind(form); err != nil {
		replyBadRequest(c, "Some fields is not valid", nil)
		return
	}

	a.withUserByID(c, func(user *models.User) {
		var resetCodeKey string
		if form.OldPassword != "" {
			if !password.CompareHashPassword(user.Password, form.OldPassword) {
				replyBadRequest(c, "Old password incorrect", nil)
				return
			}
		} else if form.ResetCode != "" {
			resetCodeKey = a.Config.Services.Account.PasswordResetCodePrefix + user.Email + "." + form.ResetCode
			_, err := a.Cache.Get(resetCodeKey)
			if err != nil {
				if err == cache.ErrNotFound {
					replyNotFound(c, "Reset code is invalid", nil)
					return
				}
				replyInternalError(c, err)
				return
			}
		} else {
			replyBadRequest(c, "Old password or reset code is missing", nil)
			return
		}

		user.Password = password.GenHashPassword(form.NewPassword)
		if err := a.ADB.UpdateUser(c.Request.Context(), user); err != nil {
			replyInternalError(c, err)
			return
		}

		if resetCodeKey != "" {
			if err := a.Cache.Del(resetCodeKey); err != nil {
				log.Logger(c.Request.Context()).Warn("Failed to delete password reset code",
					zap.String("key", resetCodeKey),
					zap.Error(err))
			}
		}
		replyOK(c, nil)
	})
}

type passwordResetCodeReply struct {
	ID string `json:"id"`
}

//CreatePasswordResetCode handles POST /user/password/reset_code
func (a *Account) CreatePasswordResetCode(c *gin.Context) {
	var form forms.PasswordResetCodeForm

	if err := c.ShouldBind(&form); err != nil {
		replyBadRequest(c, "Some fields is not valid", err)
		return
	}

	user, err := a.ADB.GetUserByEmail(context.Background(), form.Email)
	if err != nil {
		if err == db.ErrNotFound {
			replyNotFound(c, "The email does not exists", nil)
		} else {
			replyInternalError(c, err)
		}
		return
	}

	resetCode := utils.RandomDigit(6)
	if err := a.Cache.Set(&cache.Entity{
		Key:        a.Config.Services.Account.PasswordResetCodePrefix + form.Email + "." + resetCode,
		Value:      []byte(""),
		Expiration: a.Config.Services.Account.PasswordResetCodeExpiration,
	}); err != nil {
		replyInternalError(c, err)
		return
	}

	replyEmail(c, mailer.User{
		Address: a.Config.Mailer.Username,
	}, mailer.User{
		Address: user.Email,
	}, a.Config.Mailer.EmailTemplates.PasswordResetCodeName,
		map[string]string{
			"code": resetCode,
		},
		&passwordResetCodeReply{
			ID: user.ID,
		},
	)

	// emailContent := template.GenEmailContent(
	// 	a.Config.Mailer.EmailTemplates.PasswordResetCodeName,
	// 	map[string]string{
	// 		"code": resetCode,
	// 	},
	// )

	// if err := mailer.Send(&mailer.Message{
	// 	From: mailer.User{
	// 		Address: a.Config.Mailer.Username,
	// 	},
	// 	To: []mailer.User{
	// 		mailer.User{
	// 			Address: user.Email,
	// 		},
	// 	},
	// 	Subject:     emailContent.Subject,
	// 	ContentType: "text/html",
	// 	Body:        emailContent.Body,
	// }); err != nil {
	// 	log.Logger(c.Request.Context()).Warn("Failed to send email",
	// 		zap.String("from", a.Config.Mailer.Username),
	// 		zap.String("to", user.Email), zap.Error(err),
	// 	)
	// 	replyInternalError(c, err)
	// 	return
	// }

	// replyOK(c, &passwordResetCodeReply{
	// 	ID: user.ID,
	// })
}

func (a *Account) withUserByContext(c *gin.Context, f func(user *models.User)) {
	var user *models.User

	v, ok := c.Get(a.AuthMiddleware.Options().IdentityKey)
	if !ok {
		replyUnauthorized(c, "You don't have permission to access", nil)
		return
	}

	if user, ok = v.(*models.User); !ok {
		replyUnauthorized(c, "You don't have permission to access", nil)
		return
	}

	f(user)
}

func (a *Account) withUserByID(c *gin.Context, f func(user *models.User)) {
	id := c.Param("id")
	user, err := a.ADB.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if err == db.ErrNotFound {
			replyUnauthorized(c, "You don't have permission to access", nil)
			return
		}
		replyInternalError(c, err)
		return
	}

	f(user)
}

func (a *Account) userWithRoles(user *models.User) error {
	roles, err := a.RoleMgr.GetRoleForUser(user.ID)
	if err != nil {
		return err
	}
	user.Roles = strings.Join(roles, ",")
	return nil
}
