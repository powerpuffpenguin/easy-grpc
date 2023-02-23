package register

import (
	"server/logger"
	m_logger "server/m/server/logger"
	m_session "server/m/server/session"
	m_system "server/m/server/system"
	m_user "server/m/server/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GRPC(srv *grpc.Server, gateway *runtime.ServeMux, cc *grpc.ClientConn) {
	ms := []Module{
		m_system.Module(0),
		m_session.Module(0),
		m_user.Module(0),
		m_logger.Module(0),
	}
	for _, m := range ms {
		m.RegisterGRPC(srv)
		if gateway != nil {
			e := m.RegisterGateway(gateway, cc)
			if e != nil {
				logger.Logger.Panic(`register gateway error`,
					zap.Error(e),
				)
			}
		}
	}
}

type Module interface {
	RegisterGRPC(srv *grpc.Server)
	RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error
}
