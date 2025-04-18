package main

import (
	"easy-chat/apps/im/api/internal/config"
	"easy-chat/apps/im/api/internal/handler"
	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/pkg/resultx"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"sync"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	/*err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "10.0.0.17:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-api.yaml",
		ConfigFilePath: "etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

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
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}