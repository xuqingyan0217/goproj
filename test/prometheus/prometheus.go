package main

import (
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/prometheus"
	"time"
)

/*func main() {
	// 1. go语言运行的相关数据指标
	// 2. 自己自定义采集的信息
	// 创建一个新的Gauge型指标，用于跟踪当前的瞬时值。
	// Gauge适用于度量那些随时间变化的值，比如温度或当前的用户数量。
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
	    Name: "test_gauge",
	    Help: "test gauge",
	})
	// 必须注册Gauge到prometheus的默认登记处，以便Prometheus可以抓取它。
	prometheus.MustRegister(gauge)

	// 初始化变量i，用于演示目的。
	var i int
	// 启动一个新的goroutine，用于在后台不断更新Gauge的值。
	go func() {
	    for{
	        i++
	        // 当i为偶数时，增加Gauge的值。
	        // 这是一个简单的逻辑，用于模拟实际应用中可能的指标更新逻辑。
	        if i%2 == 0{
	            gauge.Inc()
	        }
	        // 暂停执行1秒，模拟指标更新的间隔。
	        time.Sleep(time.Second)
	    }
	}()

	// 配置HTTP服务器以提供Prometheus抓取的端点。
	// Prometheus通过HTTP GET请求来这个端点抓取指标数据。
	http.Handle("/metrics", promhttp.Handler())
	// 监听指定的端口，等待Prometheus来抓取指标数据。
	http.ListenAndServe(":1234", nil)
}*/

func main() {

	prometheus.StartAgent(prometheus.Config{
		Host: "0.0.0.0",
		Port: 1234,
		Path: "/metrics",
	})

	gauge := metric.NewGaugeVec(&metric.GaugeVecOpts{
		Name: "test_go_zero_gauge",
		Help: "test go-zero gauge",
	})

	// 初始化变量i，用于演示目的。
	var i int

	for{
		i++
		// 当i为偶数时，增加Gauge的值。
		// 这是一个简单的逻辑，用于模拟实际应用中可能的指标更新逻辑。
		if i%2 == 0{
			gauge.Inc()
		}
		// 暂停执行1秒，模拟指标更新的间隔。
		time.Sleep(time.Second)
	}

}
