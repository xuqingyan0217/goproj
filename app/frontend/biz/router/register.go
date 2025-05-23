// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	aieino "gomall/app/frontend/biz/router/aieino"
	auth "gomall/app/frontend/biz/router/auth"
	cart "gomall/app/frontend/biz/router/cart"
	category "gomall/app/frontend/biz/router/category"
	checkout "gomall/app/frontend/biz/router/checkout"

	home "gomall/app/frontend/biz/router/home"
	order "gomall/app/frontend/biz/router/order"
	product "gomall/app/frontend/biz/router/product"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	aieino.Register(r)

	order.Register(r)

	checkout.Register(r)

	cart.Register(r)

	category.Register(r)

	product.Register(r)

	auth.Register(r)

	home.Register(r)
}
