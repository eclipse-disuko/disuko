package rest

import (
	"context"

	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/middlewareDisco"
)

func extractPATUser(ctx context.Context) *user.User {
	raw := ctx.Value(middlewareDisco.PATUserKey)
	u, ok := raw.(*user.User)
	if !ok {
		exception.ThrowExceptionSendDeniedResponse()
	}
	return u
}
