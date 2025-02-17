package service

import (
	"context"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/checkout/infra/mq"
	"github.com/Blue-Berrys/Tiktok_e_commerce/app/checkout/infra/rpc"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/cart"
	checkout "github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/checkout"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/email"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/order"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/payment"
	"github.com/Blue-Berrys/Tiktok_e_commerce/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

// Run create note info
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	if rpc.CartClient == nil {
		rpc.InitClient()
	}

	// 1.1.get user's cart product
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
	}
	if cartResult == nil || cartResult.Cart.Items == nil {
		return nil, kerrors.NewGRPCBizStatusError(5004001, "cart is empty")
	}
	klog.Infof("查询到的购物车结果为：%v", cartResult)

	// 1.2.calculate order amount
	var (
		total float32
		oi    []*order.OrderItem
	)
	for _, cartItem := range cartResult.Cart.Items {
		productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: cartItem.ProductId,
		})

		if err != nil {
			return nil, err
		}

		if productResp.Product == nil {
			continue
		}

		p := productResp.Product.Price

		cost := p * float32(cartItem.Quantity)
		total += cost

		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: cartItem.ProductId,
				Quantity:  cartItem.Quantity,
			},
			Cost: cost,
		})
	}

	//2.create order
	var orderId string
	orderResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		Email: req.Email,
		Items: oi,
	})
	if err != nil {
		klog.Errorf("创建订单err:%v", err)
		return nil, kerrors.NewGRPCBizStatusError(5004002, err.Error())
	}
	if orderResp != nil && orderResp.Order != nil {
		orderId = orderResp.Order.OrderId
	}

	//2.2.构造支付请求
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	}

	//2.3 清空用户购物车
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		klog.Errorf("清空购物车err:%v", err)
	}

	//3.1 charge the order.结算订单
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		klog.Errorf("结算err:%v", err)
		return nil, err
	}

	//3.2订单结算的结果
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "3531095171@qq.com",
		To:          "2098245863@qq.com",
		ContentType: "text/plain",
		Subject:     "You have created an order in the Tiktok_e_commerce",
		Content:     "You have created an order in the Tiktok_e_commerce",
	})
	msg := &nats.Msg{
		Subject: "email",
		Data:    data,
		Header:  make(nats.Header),
	}
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))

	_ = mq.Nc.PublishMsg(msg)

	klog.Infof("结算订单的结果为:%v", paymentResult)

	//3.3 将订单置为paid状态
	//userId := frontutils.GetUserIdFromCtx(s.ctx)
	//_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{
	//	UserId:  uint32(userId),
	//	OrderId: orderResp.Order.OrderId,
	//})
	//if err != nil {
	//	return nil, err
	//}

	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return
}
