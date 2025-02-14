package service

import (
	"context"
	"errors"
	frontUtils "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/utils"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal/model"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/order/biz/dal/mysql"
	common "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/common"
	order "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order"
)

type UpdateOrderInfoService struct {
	ctx context.Context
} // NewUpdateOrderInfoService new UpdateOrderInfoService
func NewUpdateOrderInfoService(ctx context.Context) *UpdateOrderInfoService {
	return &UpdateOrderInfoService{ctx: ctx}
}

// Run create note info
func (s *UpdateOrderInfoService) Run(req *order.UpdateOrderInfoReq) (resp *common.Empty, err error) {
	userId := frontUtils.GetUserIdFromCtx(s.ctx)
	//基础校验
	if req.OrderId == "" {
		return nil, errors.New("userId or orderId is nil")
	}
	//检查要修改的订单所属用户是否与当前用户一致
	o, err := model.GetByOrderId(s.ctx, mysql.DB, req.OrderId)
	if err != nil {
		return nil, err
	}
	if o.UserId != uint32(userId) {
		return nil, errors.New("非法操作")
	}
	//更新订单的收货地址信息
	newAddr := req.Address
	if newAddr.StreetAddress != "" {
		o.Consignee.StreetAddress = newAddr.StreetAddress
	}
	if newAddr.ZipCode != "" {
		o.Consignee.ZipCode = newAddr.ZipCode
	}
	if newAddr.Country != "" {
		o.Consignee.Country = newAddr.Country
	}
	if newAddr.State != "" {
		o.Consignee.State = newAddr.State
	}
	if newAddr.City != "" {
		o.Consignee.City = newAddr.City
	}
	//写回
	err = model.UpdateOrderInfo(s.ctx, mysql.DB, userId, req.OrderId, o)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
