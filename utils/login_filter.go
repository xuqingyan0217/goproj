package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func CmsLoginFilter(ctx *context.Context)  {
	//获取到当前的session
	cmsUserName := ctx.Input.Session("cms_user_name")
	if cmsUserName == nil {
		//没有session，说明没有登录，让页面返回到登录界面
		ctx.Redirect(302, beego.URLFor("LoginController.Get"))
	}
}
