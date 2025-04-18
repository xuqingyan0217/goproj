package mtl

import (
	"github.com/joho/godotenv"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"os"
)

func InitTracing(serviceName string) provider.OtelProvider {
	_ = godotenv.Load()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		provider.WithInsecure(),
		// 关闭open-telemetry的指标功能，因为我们使用了Prometheus的自定义指标功能
		provider.WithEnableMetrics(false),
	)
	return p
}
