package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apierr "github.com/ngs24313/gopu/api/error"
	"github.com/ngs24313/gopu/utils/log"
	"github.com/ngs24313/gopu/utils/mailer"
	"github.com/ngs24313/gopu/utils/mailer/template"
	"go.uber.org/zap"
)

func replyError(c *gin.Context, err *apierr.AppError) {
	c.JSON(err.Code, err)
}

func replyNotFound(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusNotFound, apierr.NewAppError(
		http.StatusNotFound,
		msg,
		err,
	))
}

func replyUnauthorized(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, apierr.NewAppError(
		http.StatusUnauthorized,
		msg,
		err,
	))
}

func replyBadRequest(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusBadRequest, apierr.NewAppError(
		http.StatusBadRequest,
		msg,
		err,
	))
}

func replyInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, apierr.NewAppError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		err,
	))
}

func replyOK(c *gin.Context, data interface{}) {
	if data == nil {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, data)
}

func replyEmail(c *gin.Context,
	from mailer.User,
	to mailer.User,
	tmplName string,
	tmplData interface{},
	okData interface{}) {
	emailContent := template.GenEmailContent(tmplName, tmplData)

	if err := mailer.Send(&mailer.Message{
		From: from,
		To: []mailer.User{
			to,
		},
		Subject:     emailContent.Subject,
		ContentType: "text/html",
		Body:        emailContent.Body,
	}); err != nil {
		log.Logger(c.Request.Context()).Warn("Failed to send email",
			zap.String("from", from.Address),
			zap.String("to", to.Address), zap.Error(err),
		)
		replyInternalError(c, err)
		return
	}
	replyOK(c, okData)
}
