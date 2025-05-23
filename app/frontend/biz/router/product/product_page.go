// Code generated by hertz generator. DO NOT EDIT.

package product

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	product "gomall/app/frontend/biz/handler/product"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	root.GET("/product", append(_getproductMw(), product.GetProduct)...)
	_product := root.Group("/product", _productMw()...)
	{
		_api := _product.Group("/api", _apiMw()...)
		_api.POST("/create", append(_createproductMw(), product.CreateProduct)...)
		_api.DELETE("/delete", append(_deleteproductMw(), product.DeleteProduct)...)
		_api.POST("/update", append(_updateMw(), product.Update)...)
	}
	root.GET("/search", append(_searchproductsMw(), product.SearchProducts)...)
}
