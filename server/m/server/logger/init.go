package logger

import (
	"context"
	"server/m/server/logger/internal/db"
	grpc_logger "server/protocol/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_logger.RegisterLoggerServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	db.Init()
	return grpc_logger.RegisterLoggerHandler(context.Background(), gateway, cc)
}
