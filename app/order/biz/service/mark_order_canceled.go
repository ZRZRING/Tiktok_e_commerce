package service

import (
	"context"
	"errors"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal/model"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal/mysql"
	common "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/common"
	order "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/klog"
)

type MarkOrderCanceledService struct {
	ctx context.Context
} // NewMarkOrderCanceledService new MarkOrderCanceledService
func NewMarkOrderCanceledService(ctx context.Context) *MarkOrderCanceledService {
	return &MarkOrderCanceledService{ctx: ctx}
}

// Run create note info
func (s *MarkOrderCanceledService) Run(req *order.CancelOrderReq) (resp *common.Empty, err error) {
	//基础校验
	if req.UserId == 0 || req.OrderId == "" {
		return nil, errors.New("userId or orderId is nil")
	}
	//修改订单为canceled状态
	_, err = model.GetOrder(s.ctx, mysql.DB, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf("model.GetOrder.err:%v", err)
		return nil, err
	}
	err = model.UpdateOrderState(s.ctx, mysql.DB, req.UserId, req.OrderId, model.OrderStateCanceled)
	if err != nil {
		klog.Errorf("model.ListOrder.err:%v", err)
		return nil, err
	}
	resp = &common.Empty{}
	return
}
