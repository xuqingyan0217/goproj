package main

import (
	_ "Content/Db"     //连接数据库
	_ "Content/models" //创建模型

	_ "Content/routers" //寻找路由

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

func main() {
	//添加日志模块，采用所有级别的日志信息都放在一处
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/test.log"}`)
	//定义过滤器函数，用于配合session使用，注意第三个参数是函数名，不要加括号
	//beego.InsertFilter("/cms/main/*", beego.BeforeRouter, utils.CmsLoginFilter)
	//beego.InsertFilter("/front/main/*", beego.BeforeRouter, utils.CmsLoginFilter)
	//自定义模板函数
	beego.AddFuncMap("ShowPrePage", HandlePrePage)
	beego.AddFuncMap("ShowNextPage", HandleNextPage)
	beego.SetStaticPath("/static", "static")
	beego.Run()
}

func HandlePrePage(data int) string {
	pageIndex := data - 1
	PageIndex := strconv.Itoa(pageIndex)
	return PageIndex
}
func HandleNextPage(data int) int {
	pageIndex := data + 1
	return pageIndex
}
