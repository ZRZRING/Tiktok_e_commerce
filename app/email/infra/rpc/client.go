package rpc

import (
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/email/conf"
	frontUtils "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/utils"
	"github.com/Blue-Berrys/Tiktok_e_commerce/common/clientsuite"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/kitex/client"
	"sync"
)

var (
	OrderClient orderservice.Client

	once         sync.Once
	err          error
	serviceName  string
	metricsPort  string
	registryAddr string
)

func Init() {
	once.Do(func() {
		serviceName = conf.GetConf().Kitex.Service
		metricsPort = conf.GetConf().Kitex.MetricsPort
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		initOrderClient()
	})
}

func initOrderClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: serviceName,
			RegistryAddr:       registryAddr,
		}),
	}

	OrderClient, err = orderservice.NewClient("order", opts...)
	frontUtils.MustHandleError(err)
}
