package service

import (
	"context"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal"
	order "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMarkOrderCanceled_Run(t *testing.T) {
	ctx := context.Background()
	s := NewMarkOrderCanceledService(ctx)

	//chdir to adjust config file path
	if err := os.Chdir("../.."); err != nil {
		log.Fatalf("chdir err : %v", err)
	}

	//init database
	_ = godotenv.Load()
	dal.Init()

	req := &order.CancelOrderReq{
		OrderId: "c5d8ec32-7344-4b69-b238-a291dbcd1a85",
		UserId:  2,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
