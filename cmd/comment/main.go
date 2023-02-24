package main

import (
	"TikV/dal"
	comment "TikV/kitex_gen/comment/commentsrv"
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
	Config      = ttviper.ConfigInit("TIKV_COMMENT", "commentConfig")
	ServiceName = Config.Viper.GetString("Server.Name")
	ServiceAddr = fmt.Sprintf("%s:%d", Config.Viper.GetString("Server.Address"), Config.Viper.GetInt("Server.Port"))
	EtcdAddress = fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
	Jwt         *jwt.JWT
)

// Comment RPC Server 端配置初始化
func Init() {
	dal.Init()
	Jwt = jwt.NewJWT([]byte(Config.Viper.GetString("JWT.signingKey")))
}

// Comment RPC Server 端运行
func main() {
	var logger = dlog.InitLog(3)
	defer logger.Sync()

	klog.SetLogger(logger)

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

	svr := comment.NewServer(
		new(CommentSrvImpl),
		server.WithServiceAddr(addr),                                       // address
		server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}),
	)

	if err := svr.Run(); err != nil {
		klog.Fatalf("%s stopped with error:", ServiceName, err)
	}
}
