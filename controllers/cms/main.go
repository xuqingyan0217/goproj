package cms

import "github.com/astaxie/beego"

type MainController struct {
	beego.Controller
}

func (m *MainController)Get ()  {
	Username := m.GetSession("cms_user_name")
	flag := false
	if Username != nil {
		flag = true
	}
	m.Data["Flag"] = flag
	m.Data["Username"] = Username
	m.TplName = "cms/index.html"
}

func (m *MainController) Welcome () {
	m.TplName = "cms/welcome.html"
}