package rpc

import (
	"common/clientsuite"
	"github.com/cloudwego/kitex/client"
	"gomall/app/cart/conf"

	cartUtils "gomall/app/cart/utils"
	"gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"sync"
)

var (
	ProductClient productcatalogservice.Client
	// 确保只初始化一次
	once         sync.Once
	err          error
	ServiceName  = conf.GetConf().Kitex.Service
	RegistryAddr = conf.GetConf().Registry.RegistryAddress
)

func Init() {
	once.Do(func() {
		initProductClient()
	})
}
func initProductClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}

	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	cartUtils.MustHandlerError(err)
}
