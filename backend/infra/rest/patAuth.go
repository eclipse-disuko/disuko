package rest

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/golang-jwt/jwt/v4"
)

func patAuth(rs *logy.RequestSession, pr *project.Project, prRepo projectRepo.IProjectRepository, userRepo userRepo.IUsersRepository, tokenStr string) string {
	token, err := jwt.ParseWithClaims(tokenStr, &user.UserTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.Config.Auth.UserTokenSigningKey), nil
	})
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), err.Error())
	}
	claims, ok := token.Claims.(*user.UserTokenClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Unexpected claims")
	}
	user := userRepo.FindByKey(rs, claims.UserKey, false)
	if user == nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Unexpected claims")
	}
	ut := user.Token(claims.TokenKey)
	if ut == nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Unexpected claims")
	}
	if ut.Expired() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Unexpected claims")
	}
	if m := pr.GetMember(user.User); m != nil && m.UserType == project.OWNER {
		return user.TokenOrigin(ut)
	}
	exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Project access denied")
	return ""
}
