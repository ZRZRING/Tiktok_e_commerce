package service

import (
	"context"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	order "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order"
)

func TestUpdateOrderInfo_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateOrderInfoService(ctx)

	//chdir to adjust config file path
	if err := os.Chdir("../.."); err != nil {
		log.Fatalf("chdir err : %v", err)
	}

	//init database
	_ = godotenv.Load()
	dal.Init()

	//todo:这个测试用例没有身份信息
	req := &order.UpdateOrderInfoReq{
		OrderId: "c5d8ec32-7344-4b69-b238-a291dbcd1a85",
		Address: &order.Address{
			StreetAddress: "",
			City:          "guangzhou",
			State:         "",
			Country:       "",
			ZipCode:       "",
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
