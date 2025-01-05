package cms

import "github.com/astaxie/beego"

type LogoutController struct {
	beego.Controller
}

func (l *LogoutController) Logout()  {
	l.DelSession("cms_user_name")
	l.Redirect(beego.URLFor("MainController.Get"), 302)
}
