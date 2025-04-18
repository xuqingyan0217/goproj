module my_assistant

go 1.23.6

require (
	github.com/cloudwego/eino v0.3.23
	github.com/cloudwego/eino-ext/callbacks/apmplus v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/callbacks/langfuse v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/document/loader/file v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/embedding/ark v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/indexer/redis v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/model/ark v0.1.6
	github.com/cloudwego/eino-ext/components/retriever/redis v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/tool/duckduckgo v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/components/tool/mcp v0.0.0-20250416093724-5f53a5d6cc28
	github.com/cloudwego/eino-ext/devops v0.1.7
	github.com/cloudwego/hertz v0.9.7
	github.com/go-sql-driver/mysql v1.9.2
	github.com/google/uuid v1.6.0
	github.com/hertz-contrib/obs-opentelemetry/provider v0.3.0
	github.com/hertz-contrib/obs-opentelemetry/tracing v0.4.1
	github.com/hertz-contrib/sse v0.1.0
	github.com/joho/godotenv v1.5.1
	github.com/mark3labs/mcp-go v0.20.1
	github.com/redis/go-redis/v9 v9.7.3
	github.com/xuri/excelize/v2 v2.9.0
	go.opentelemetry.io/otel v1.35.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bytedance/gopkg v0.1.0 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/eino-ext/libs/acl/langfuse v0.0.0-20250409060521-ba8646352e4b // indirect
	github.com/cloudwego/eino-ext/libs/acl/opentelemetry v0.0.0-20250225080340-5935633151d3 // indirect
	github.com/cloudwego/netpoll v0.6.4 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/getkin/kin-openapi v0.118.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1 // indirect
	github.com/invopop/yaml v0.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matoous/go-nanoid v1.5.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/perimeterx/marshmallow v1.1.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slongfield/pyfmt v0.0.0-20220222012616-ea85ff4c361f // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/volcengine/volc-sdk-golang v1.0.23 // indirect
	github.com/volcengine/volcengine-go-sdk v1.0.185 // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.59.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.20.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.20.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.34.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.34.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
