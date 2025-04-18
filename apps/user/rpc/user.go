package main

import (
	"easy-chat/apps/user/rpc/internal/config"
	"easy-chat/apps/user/rpc/internal/server"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/interceptor/rpcserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

var wg sync.WaitGroup

var grpcSvr = grpc.NewServer()

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	/*err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "10.0.0.17:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "user",
		Configs:        "user-rpc.yaml",
		ConfigFilePath: "etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅退出
		// proc.WrapUp()
		// 确保在重新加载配置时，旧的gRPC服务器已经停止

		// 优雅退出
		proc.WrapUp()

		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(any(err))
	}
	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)
	wg.Wait()*/
	Run(c)

}

func Run(c config.Config) {
	ctx := svc.NewServiceContext(c)

	if err := ctx.SetRootToken(); err != nil {
		panic(any(err))
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
