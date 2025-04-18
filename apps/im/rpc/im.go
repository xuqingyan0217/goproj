package main

import (
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/server"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/interceptor/rpcserver"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/reflection"
	"sync"

	"easy-chat/apps/im/rpc/internal/config"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

var wg sync.WaitGroup

var grpcSvr *grpc.Server

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	/*err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "10.0.0.17:3379",
		ProjectKey:     "c22c60349630d688cef20a3fd708ad87",
		Namespace:      "im",
		Configs:        "im-rpc.yaml",
		ConfigFilePath: "etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅退出
		// proc.WrapUp()
		grpcSvr.GracefulStop()

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

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		im.RegisterImServer(grpcServer, server.NewImServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}