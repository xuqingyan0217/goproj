package rpc

import (
	"common/clientsuite"
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/config-consul/consul"
	"gomall/rpc_gen/kitex_gen/AIEino/aieinoservice"

	consulclient "github.com/kitex-contrib/config-consul/client"
	"gomall/app/frontend/conf"
	frontendUtils "gomall/app/frontend/utils"
	"gomall/rpc_gen/kitex_gen/cart/cartservice"
	"gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"gomall/rpc_gen/kitex_gen/order/orderservice"
	"gomall/rpc_gen/kitex_gen/product"
	"gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"gomall/rpc_gen/kitex_gen/user/userservice"
	"sync"
)

var (
	UserClient     userservice.Client
	ProductClient  productcatalogservice.Client
	CartClient     cartservice.Client
	CheckoutClient checkoutservice.Client
	OrderClient    orderservice.Client
	AIEinoClient   aieinoservice.Client
	// 确保只初始化一次
	once         sync.Once
	ServiceName  = frontendUtils.ServiceName
	RegistryAddr = conf.GetConf().Hertz.RegistryAddr
	err          error
)

func Init() {
	once.Do(func() {
		initUserClient()
		initProductClient()
		initCartClient()
		initCheckoutClient()
		initOrderClient()
		initAIEinoClient()
	})
}
func initUserClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	UserClient, err = userservice.NewClient("user", opts...)
	frontendUtils.HandlerError(err)
}

func initProductClient() {
	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	cbs.UpdateServiceCBConfig("frontend/product/GetProduct",
		circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 2},
	)
	consulClient, err := consul.NewClient(consul.Options{
		Addr: "consul-svc:8500",
	})
	frontendUtils.HandlerError(err)

	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
		client.WithCircuitBreaker(cbs),
		client.WithFallback(fallback.NewFallbackPolicy(
			fallback.UnwrapHelper(
				func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
					if err == nil {
						return resp, nil
					}
					methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
					if methodName != "ListProducts" {
						return resp, err
					}
					return &product.ListProductsResp{
						Products: []*product.Product{
							{
								Price:       6.6,
								Id:          3,
								Picture:     "/static/image/t-shirt.jpeg",
								Name:        "T-shirt",
								Description: "CloudWeGo shirt",
							},
						},
					}, nil
				},
			),
		)),
		client.WithSuite(consulclient.NewSuite("product", ServiceName, consulClient)),
	}
	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	frontendUtils.HandlerError(err)
}

func initCartClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	CartClient, err = cartservice.NewClient("cart", opts...)
	frontendUtils.HandlerError(err)
}

func initCheckoutClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	CheckoutClient, err = checkoutservice.NewClient("checkout", opts...)
	frontendUtils.HandlerError(err)
}

func initOrderClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	OrderClient, err = orderservice.NewClient("order", opts...)
	frontendUtils.HandlerError(err)
}

func initAIEinoClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	AIEinoClient, err = aieinoservice.NewClient("aieino", opts...)
	frontendUtils.HandlerError(err)
}
