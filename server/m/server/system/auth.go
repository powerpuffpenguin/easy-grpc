package system

import (
	"context"
	"server/db"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	ctx, session, e := s.Session(ctx)
	if e != nil {
		return ctx, e
	} else if session.AuthAny(db.Root) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
