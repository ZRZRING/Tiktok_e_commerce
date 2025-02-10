// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	auth "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/auth"
	cart "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/cart"
	category "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/category"
	checkout "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/checkout"
	home "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/home"
	product "github.com/Blue-Berrys/Tiktok_e_commerce/app/frontend/biz/router/product"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	checkout.Register(r)

	cart.Register(r)

	category.Register(r)

	product.Register(r)

	auth.Register(r)

	home.Register(r)
}
