package front

import (
	"Content/models"
	"Content/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"regexp"
)

type RegisterController struct {
	beego.Controller
}

func (r *RegisterController) Get() {
	r.TplName = "front/register.html"
}

// Post 处理用户注册的POST请求
func (r *RegisterController) Post() {
	// 获取前端传来的用户名、密码和确认密码
	username := r.GetString("username")
	password := r.GetString("password")
	repwd := r.GetString("repassword")

	// 后端验证用户名格式（参考前端验证逻辑，确保更严谨的验证）
	if ok := validateUsername(username); !ok {
		r.Data["json"] = map[string]interface{}{"error": "用户名格式不正确，请输入3-20位包含字母、数字、下划线的字符"}
		r.ServeJSON()
		return
	}

	// 后端验证密码格式（同样参考前端验证思路）
	if ok := validatePassword(password); !ok {
		r.Data["json"] = map[string]interface{}{"error": "密码格式不正确，请输入6-20位包含字母、数字、常见特殊字符的字符"}
		r.ServeJSON()
		return
	}

	// 后端再次验证密码一致性（前端验证可能被绕过，多一层后端验证更安全）
	if password != repwd {
		r.Data["json"] = map[string]interface{}{"error": "两次的密码不一致"}
		r.ServeJSON()
		return
	}

	// 密码加密处理
	md5_password := utils.GetMd5(password)

	// 创建数据库操作对象（假设使用了某种ORM框架，此处以简单示意）
	o := orm.NewOrm()
	user := models.User{
		UserName: username,
		Password: md5_password,
		Cover:    "static/upload/bq3.png",
	}

	// 插入用户信息到数据库
	_, err := o.Insert(&user)
	if err != nil {
		r.Data["json"] = map[string]interface{}{"error": "注册失败，数据库插入错误"}
		r.ServeJSON()
		return
	}

	// 注册成功，返回成功提示信息（可根据前端实际需求调整具体返回数据结构）
	r.Data["json"] = map[string]interface{}{"message": "注册成功"}
	r.ServeJSON()
}

// 验证用户名格式的函数（参考前端验证逻辑，保持一致性）
func validateUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_]{3,20}$`
	return len(username) >= 3 && len(username) <= 20 && matchRegex(pattern, username)
}

// 验证密码格式的函数（参考前端验证逻辑，保持一致性）
func validatePassword(password string) bool {
	pattern := `^[a-zA-Z0-9@#$%^&+=]{6,20}$`
	return len(password) >= 6 && len(password) <= 20 && matchRegex(pattern, password)
}

// 辅助函数，用于判断字符串是否匹配给定的正则表达式
func matchRegex(pattern, str string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}
