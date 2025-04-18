package routers

import (
	"Content/controllers/cms"
	"Content/controllers/front"
	"github.com/astaxie/beego"
)

func init() {
	//后端
	//登录界面
	beego.Router("/cms", &cms.LoginController{})
	//退出登录
	beego.Router("/cms/logout", &cms.LogoutController{}, "get:Logout")
	//登录后主界面
	beego.Router("/", &cms.MainController{})
	//主界面下面的子界面
	beego.Router("/cms/welcome", &cms.MainController{}, "get:Welcome")
	//主界面下的日记列表
	beego.Router("/cms/post_list", &cms.PostController{}, "get:Get")
	//主界面下添加日记
	beego.Router("/cms/main/post_add", &cms.PostController{}, "get:ToAdd;post:DoAdd")
	//主界面下的删除日记
	beego.Router("/cms/main/post_delete", &cms.PostController{}, "get:PostDelete")
	//主界面下编辑日记
	beego.Router("/cms/main/post_edit", &cms.PostController{}, "get:ToEdit;post:DoEdit")

	//前端
	//个人页列表
	beego.Router("/front/main", &front.IndexController{})
	//列表详情页
	beego.Router("/front/main/detail", &front.IndexController{}, "get:PostDetail")
	//详情页下的评论
	beego.Router("/front/main/comment", &front.CommentController{})
	// 注册页
	beego.Router("/register", &front.RegisterController{})
}
