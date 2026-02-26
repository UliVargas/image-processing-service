package auth

import "context"

type AuthenticatedUser struct {
	UserID string
	JTI    string
}

type contextKey string

const AuthKey contextKey = "auth_info"

func GetAuthUser(ctx context.Context) (AuthenticatedUser, bool) {
	user, ok := ctx.Value(AuthKey).(AuthenticatedUser)
	return user, ok
}
