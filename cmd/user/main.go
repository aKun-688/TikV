package main

import (
	"TikV/cmd/user/service"
	"TikV/dal"
	user "TikV/kitex_gen/user/usersrv"
	"TikV/pkg/dlog"
	"TikV/pkg/jwt"
	middleware "TikV/pkg/mw"
	"TikV/pkg/ttviper"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"net"
)

var (
	Config       = ttviper.ConfigInit("TIKV_USER", "userConfig")
	ServiceName  = Config.Viper.GetString("Server.Name")
	ServiceAddr  = fmt.Sprintf("%s:%d", Config.Viper.GetString("Server.Address"), Config.Viper.GetInt("Server.Port"))
	EtcdAddress  = fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
	Jwt          *jwt.JWT
	Argon2Config *service.Argon2Params
)

// User RPC Server 端配置初始化
func Init() {
	dal.Init()
	Jwt = jwt.NewJWT([]byte(Config.Viper.GetString("JWT.signingKey")))
	Argon2Config = &service.Argon2Params{
		Memory:      Config.Viper.GetUint32("Server.Argon2ID.Memory"),
		Iterations:  Config.Viper.GetUint32("Server.Argon2ID.Iterations"),
		Parallelism: uint8(Config.Viper.GetUint("Server.Argon2ID.Parallelism")),
		SaltLength:  Config.Viper.GetUint32("Server.Argon2ID.SaltLength"),
		KeyLength:   Config.Viper.GetUint32("Server.Argon2ID.KeyLength"),
	}
}

func main() {
	var logger = dlog.InitLog(3)
	defer logger.Sync()

	klog.SetLogger(logger)

	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{EtcdAddress})
	if err != nil {
		klog.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", ServiceAddr)
	if err != nil {
		klog.Fatal(err)
	}

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	Init()

	svr := user.NewServer(new(UserSrvImpl),
		server.WithServiceAddr(addr),                                       // address
		server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}))

	if err := svr.Run(); err != nil {
		klog.Fatalf("%s stopped with error:", ServiceName, err)
	}
}
