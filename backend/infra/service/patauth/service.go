package patauth

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	Repo userRepo.IUsersRepository
}

func (s *Service) ValidateForProject(rs *logy.RequestSession, pr *project.Project, tokenStr string) string {
	user, ut := s.Validate(rs, tokenStr)
	if m := pr.GetMember(user.User); m != nil && m.UserType == project.OWNER {
		return user.TokenOrigin(ut)
	}
	exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Project access denied")
	return ""
}

func (s *Service) Validate(rs *logy.RequestSession, tokenStr string) (*user.User, *user.Token) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.UserTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.Config.PATAuth.SigningKey), nil
	})
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), err.Error())
	}
	claims, ok := token.Claims.(*user.UserTokenClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid PAT"), "Unexpected claims")
	}
	user := s.Repo.FindByKey(rs, claims.UserKey, false)
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
	return user, ut
}
