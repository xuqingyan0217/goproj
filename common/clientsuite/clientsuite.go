package clientsuite

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

type CommonClientSuite struct {
	CurrentServiceName string
	RegistryAddr       []string
}

func (s CommonClientSuite) Options() []client.Option {
	opts := []client.Option{
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		//client.WithTransportProtocol(transport.GRPC),
		client.WithShortConnection(), // 使用短链接,在linux运行时再换回去，不换也可以
		client.WithSuite(tracing.NewClientSuite()),
	}

	r, err := etcd.NewEtcdResolver(s.RegistryAddr)
	if err != nil {
		panic(err)
	}
	opts = append(opts, client.WithResolver(r))
	return opts
}
