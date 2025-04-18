/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"easy-chat/apps/im/ws/internal/handler"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/ctxdata"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"net/http"
	"sync"
	"time"

	"easy-chat/apps/im/ws/internal/config"
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
		Configs:        "im-ws.yaml",
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
	if err := c.SetUp(); err != nil {
		panic(any(err))
	}
	ctx := svc.NewServiceContext(c)
	if ctx == nil {
		panic(any("failed to initialize service context"))
	}

	// 设置服务认证的token
	token, err := ctxdata.GetJwtToken(c.JwtAuth.AccessSecret, time.Now().Unix(), 3153600000, fmt.Sprintf("ws:%s", time.Now().Unix()))
	if err != nil || token == "" {
		panic(any("failed to get JWT token"))
	}

		srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.NoAck),
		websocket.WithServerMaxConnectionIdle(20*time.Second),
		websocket.WithServerDiscover(websocket.NewRedisDiscover(http.Header{
			"Authorization": []string{token},
		}, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	)
	if srv == nil {
		panic(any("failed to initialize WebSocket server"))
	}
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
