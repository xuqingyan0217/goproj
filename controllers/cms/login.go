package cms

import (
	"Content/models"
	"Content/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (l *LoginController) Get() {
	l.TplName = "cms/login.html"
}

func (l *LoginController) Post() {
	username := l.GetString("username")
	pwd := l.GetString("password")
	// 可以在这里添加更完善的输入校验逻辑，比如判断用户名是否为空等情况
	if username == "" {
		l.Data["json"] = map[string]interface{}{
			"code":  400,
			"error": "用户名不能为空",
		}
		l.ServeJSON()
		return
	}
	if pwd == "" {
		l.Data["json"] = map[string]interface{}{
			"code":  400,
			"error": "密码不能为空",
		}
		l.ServeJSON()
		return
	}
	// 特殊字符过滤（这里假设你有合适的过滤函数，如果没有可补充完善，目前仅示意位置）
	// 例如对username进行过滤处理，此处代码省略具体过滤逻辑实现
	// username = FilterSpecialChars(username)

	// 长度校验（可以添加更严格合理的密码长度校验逻辑，目前简单示意）
	if len(pwd) < 6 {
		l.Data["json"] = map[string]interface{}{
			"code":  400,
			"error": "密码长度不能小于6位",
		}
		l.ServeJSON()
		return
	}
	// 加密
	md5_pwd := utils.GetMd5(pwd)
	// 使用orm
	o := orm.NewOrm()
	// 直接使用查询的存在语句来判断该用户是否存在
	exits := o.QueryTable(new(models.User)).Filter("user_name", username).Filter("password", md5_pwd).Exist()
	if exits {
		// 登录成功后设置session，我们的session是依据用户名设置的，所以要保证数据库了name唯一，多了cms是为了区分
		l.SetSession("cms_user_name", username)
		// 存在则返回成功信息给前端（这里返回的结构体字段可根据前端期望和实际项目需求调整）
		l.Data["json"] = map[string]interface{}{
			"code":    200,
			"message": "登录成功",
		}
		logs.Info(username, "登录成功")
	} else {
		// 不存在则返回相应错误信息给前端
		l.Data["json"] = map[string]interface{}{
			"code":  401,
			"error": "用户名或密码错误",
		}
	}
	l.ServeJSON()
}
