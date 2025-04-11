package tests

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	v1 "github.com/seth16888/wxtoken/api/v1"
	"github.com/seth16888/wxtoken/internal/di"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthsvc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

var (
	grpcServer *grpc.Server
	listener   *bufconn.Listener
	configFile string
)

func startTestGrpcServer(configFile string) (*grpc.Server, *bufconn.Listener) {
	appDeps := di.NewContainer(configFile)

	listener := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s, appDeps.Svc)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return s, listener
}

func registerServices(s *grpc.Server, svc v1.TokenServer) {
	// 注册服务
	v1.RegisterTokenServer(s, svc)
	healthpb.RegisterHealthServer(s, healthsvc.NewServer())

}

func createGrpcConnection(t *testing.T) *grpc.ClientConn {
	// 定义一个 bufconn 拨号器，用于在测试中建立连接
	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return listener.Dial()
	}

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 使用 grpc.DialContext 建立连接
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufconnDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	return conn
}

func TestMain(m *testing.M) {
	configFile = ""
	envPath := os.Getenv("WXTOKEN_CONF")
	if envPath != "" {
		configFile = envPath
	} else {
		configFile = "d:/code2025/cowx/wxtoken/conf/conf.yaml"
	}
	if configFile == "" {
		log.Fatalf("Env WXTOKEN_CONF is not set")
	}

	grpcServer, listener = startTestGrpcServer(configFile)
	time.Sleep(2 * time.Second)

	// 运行所有测试用例
	code := m.Run()

	// 清理资源
	grpcServer.GracefulStop()
	listener.Close()

	os.Exit(code)
}

// HealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	conn := createGrpcConnection(t)
	defer conn.Close()

  Convey("When service alive", t, func() {
		client := healthpb.NewHealthClient(conn)
		req := healthpb.HealthCheckRequest{
		}
		resp, err := client.Check(context.Background(), &req)

		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Status, ShouldEqual, healthpb.HealthCheckResponse_SERVING)
	})
}

// TestGetToken 测试获取token
func TestGetToken(t *testing.T) {
  conn := createGrpcConnection(t)
	defer conn.Close()

  Convey("When get token", t, func() {
    client := v1.NewTokenClient(conn)
    req := &v1.GetTokenRequest{
      AppId: 1001,
      MpId: "string",
    }

    resp, err := client.GetAccessToken(context.Background(), req)
    So(err, ShouldBeNil)
    So(resp, ShouldNotBeNil)
    So(resp.AccessToken, ShouldNotBeEmpty)
  })
}
