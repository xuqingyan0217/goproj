/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/apps/task/mq/internal/handler"
	"easy-chat/apps/task/mq/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/grpc"
	"sync"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

var wg sync.WaitGroup

var grpcSvr *grpc.Server

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	/*err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "10.0.0.17:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
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
	if err := c.SetUp(); err != nil {
		panic(any(err))
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	serviceGroup.Start()
}