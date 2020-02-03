package password

import (
	"context"

	"github.com/ngs24313/gopu/utils/log"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//GenHashPassword generate hashed password
func GenHashPassword(pwd string) string {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Logger(context.Background()).Warn("Failed to generate password by bcrypt", zap.Error(err))
	}
	return string(hashedPwd)
}

//CompareHashPassword compares a hashed password with plaintext
func CompareHashPassword(hashPwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	return err == nil
}
