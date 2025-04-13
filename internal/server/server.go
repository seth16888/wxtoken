package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/seth16888/wxtoken/api/v1"
	"github.com/seth16888/wxtoken/internal/di"
	"github.com/seth16888/wxtoken/internal/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	healthsvc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Start 启动服务
func Start(deps *di.Container) error {
	listenAddr := deps.Conf.Server.Addr
	if len(listenAddr) == 0 {
		listenAddr = ":10103"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		deps.Log.Error("failed to listen", zap.Error(err))
		return err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.TimeoutInterceptor(),
			middleware.RequestID(),
			middleware.LoggingInterceptor(deps.Log),
			middleware.ClientDisconnectInterceptor(),
			middleware.RecoverInterceptor(deps.Log),
		),
    // 明文传输
    grpc.Creds(insecure.NewCredentials()),
	)
	v1.RegisterTokenServer(s, deps.Svc)
	// 健康检查
	healthSvc := healthsvc.NewServer()
	healthpb.RegisterHealthServer(s, healthSvc)
	updateHealthStatus(healthSvc, v1.Token_ServiceDesc.ServiceName,
		healthpb.HealthCheckResponse_SERVING)

	deps.Log.Info("starting grpc server", zap.String("addr", listenAddr))
	errCh := make(chan error, 1)
	go func() {
		if err := s.Serve(listener); err != grpc.ErrServerStopped {
			deps.Log.Error("failed to serve", zap.Error(err))
			errCh <- err
		} else {
			deps.Log.Info("grpc server stopped")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		updateHealthStatus(healthSvc, v1.Token_ServiceDesc.ServiceName,
			healthpb.HealthCheckResponse_NOT_SERVING)
		deps.Log.Info("shutting down grpc server gracefully...")
		s.GracefulStop()
		deps.Log.Sync() // 确保日志同步
	case err := <-errCh:
		updateHealthStatus(healthSvc, v1.Token_ServiceDesc.ServiceName,
			healthpb.HealthCheckResponse_NOT_SERVING)
		deps.Log.Sync() // 确保日志同步
		return err
	}

  deps.Log.Info("server shutdown completed")
	return nil
}

// updateHealthStatus 更新健康检查状态
func updateHealthStatus(
	h *healthsvc.Server,
	service string,
	status healthpb.HealthCheckResponse_ServingStatus,
) {
	h.SetServingStatus(service, status)
}
