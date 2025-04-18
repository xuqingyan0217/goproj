package mtl

import (
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
)

var Registry *prometheus.Registry

func InitMetric(serviceName, metricsPort string, registryAddr []string) (registry.Registry, *registry.Info) {
	// 初始化
	Registry = prometheus.NewRegistry()
	// 注册go运行相关指标
	Registry.MustRegister(collectors.NewGoCollector())
	// 注册进程相关指标
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// 将Prometheus注册到etcd，这样就很容易知道有哪些服务了，无需做压帽
	r, _ := etcd.NewEtcdRegistry(registryAddr)
	// 转换为tcp服务地址
	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)
	// 构造注册信息
	registryInfo := &registry.Info{
		ServiceName: "prometheus",
		Addr:        addr,
		Weight:      1,
		Tags:        map[string]string{"serviceName": serviceName},
	}
	// 注册服务
	_ = r.Register(registryInfo)
	// shutdown hook，当服务关闭的时候，下线其在注册中心的信息
	server.RegisterShutdownHook(func() {
		_ = r.Deregister(registryInfo)
	})
	// 启动Prometheus服务
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	// 异步启动一个server供Prometheus拉取指标
	go http.ListenAndServe(metricsPort, nil)

	return r, registryInfo
}
