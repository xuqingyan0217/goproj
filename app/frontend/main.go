// Code generated by hertz generator.

package main

import (
	"common/mtl"
	"context"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/joho/godotenv"
	"gomall/app/frontend/infra/rpc"
	"gomall/app/frontend/middleware"
	frontentUtils "gomall/app/frontend/utils"

	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/logger/accesslog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	hertzzobslogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"


	"github.com/hertz-contrib/pprof"
	"go.uber.org/zap/zapcore"
	"gomall/app/frontend/biz/router"
	"gomall/app/frontend/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ServiceName  = frontentUtils.ServiceName
	MetricsPort  = conf.GetConf().Hertz.MetricsPort
	RegistryAddr = conf.GetConf().Hertz.RegistryAddr
)

func main() {
	_ = godotenv.Load()

	p := mtl.InitTracing(ServiceName)
	defer p.Shutdown(context.Background())

	etcd, registryInfo := mtl.InitMetric(ServiceName, MetricsPort, RegistryAddr)
	defer etcd.Deregister(registryInfo)
	// init dal
	// dal.Init()
	// 调用初始化
	rpc.Init()
	// 初始化JWT
	frontentUtils.InitJWT()
	address := conf.GetConf().Hertz.Address

	tracer, cfg := hertztracing.NewServerTracer()
	h := server.New(server.WithHostPorts(address),
		server.WithTracer(prometheus.NewServerTracer("", "", prometheus.WithDisableServer(true),
			prometheus.WithRegistry(mtl.Registry),
		)),
		tracer,
	)
	h.Use(hertztracing.ServerMiddleware(cfg))

	registerMiddleware(h)

	// add a ping route to test
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	router.GeneratedRegister(h)
	h.Spin()
}

func registerMiddleware(h *server.Hertz) {
	// 移除session相关配置
	
	// log,采用带链路跟踪的
	logger := hertzzobslogrus.NewLogger(hertzzobslogrus.WithLogger(hertzlogrus.NewLogger().Logger()))
	//loger := hertzlogrus.NewLogger()
	hlog.SetLogger(logger)
	hlog.SetLevel(conf.LogLevel())
	// 刷盘
	var flushInterval time.Duration
	if os.Getenv("GO_ENV") == "online" {
		flushInterval = time.Minute
	} else {
		flushInterval = time.Second
	}
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Hertz.LogFileName,
			MaxSize:    conf.GetConf().Hertz.LogMaxSize,
			MaxBackups: conf.GetConf().Hertz.LogMaxBackups,
			MaxAge:     conf.GetConf().Hertz.LogMaxAge,
		}),
		FlushInterval: flushInterval,
	}
	hlog.SetOutput(asyncWriter)
	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		asyncWriter.Sync()
	})

	// pprof
	if conf.GetConf().Hertz.EnablePprof {
		pprof.Register(h)
	}

	// gzip
	if conf.GetConf().Hertz.EnableGzip {
		h.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// access log
	if conf.GetConf().Hertz.EnableAccessLog {
		h.Use(accesslog.New())
	}

	// recovery
	h.Use(recovery.Recovery())

	// cores
	h.Use(cors.Default())

	middleware.Register(h)
}
