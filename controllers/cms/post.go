package cms

import (
	"Content/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	"mime/multipart"
	"strconv"
	"time"
)

type PostController struct {
	beego.Controller
}

func (p *PostController) Get() {
	o := orm.NewOrm()
	posts := []models.Post{}
	qs := o.QueryTable(new(models.Post))
	count, _ := qs.RelatedSel().Count()

	pageSize := 3
	pageIndexStr := p.GetString("pageIndex")
	logs.Info(pageIndexStr)
	PageIndex := 1
	if pageIndexStr != "" {
		var err error
		PageIndex, err = strconv.Atoi(pageIndexStr)
		if err != nil {
			logs.Error("Invalid page index:", err)
		}
	}

	FirstPage := PageIndex == 1
	EndPage := false
	start := pageSize * (PageIndex - 1)

	pageCount := math.Ceil(float64(count) / float64(pageSize))
	qs.Limit(pageSize, start).RelatedSel().All(&posts)
	if PageIndex == int(pageCount) {
		EndPage = true
	}

	Username := p.GetSession("cms_user_name")
	flag := Username != nil

	isAjax := p.Ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"

	if isAjax {
		p.Data["json"] = map[string]interface{}{
			"posts":     posts,
			"count":     count,
			"pageCount": pageCount,
			"EndPage":   EndPage,
			"FirstPage": FirstPage,
			"pageIndex": PageIndex,
			"flag":      flag,
		}
		p.ServeJSON()
		return
	}

	p.Data["Flag"] = flag
	p.Data["posts"] = posts
	p.Data["count"] = count
	p.Data["pageCount"] = pageCount
	p.Data["EndPage"] = EndPage
	p.Data["FirstPage"] = FirstPage
	p.Data["pageIndex"] = PageIndex
	p.TplName = "cms/post-list.html"
}

func (p *PostController) ToAdd() {
	p.TplName = "cms/post-add.html"
}

// DoAdd 添加文章的post请求
func (p *PostController) DoAdd() {
	title := p.GetString("title")
	desc := p.GetString("desc")
	content := p.GetString("content")
	//获取文件
	f, h, err := p.GetFile("cover")
	var cover string
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			logs.Error(err)
		}
	}(f)
	//如果cover是没上传，说明err出错，我们给个默认图片
	//如果上传了，err不为空，那么我们要把文件保存到本地然后存到数据库里面
	if err != nil {
		cover = "static/upload/user/no_pic.jpg"
	} else {
		//生成时间戳，防止重名
		timeUnix := time.Now().Unix()               //int64
		time_str := strconv.FormatInt(timeUnix, 10) //int64转换为字符串
		//保存获取到的文件
		path := "static/upload/user/" + time_str + h.Filename
		err1 := p.SaveToFile("cover", path)
		//保存出错了我们也要给个默认值
		if err1 != nil {
			cover = "static/upload/user/no_pic.jpg"
		} else {
			cover = path
		}
	}

	o := orm.NewOrm()
	//通过session获取当前用户名，然后去查询该用户并且赋值给一个user结构体
	author := p.GetSession("cms_user_name")
	user := models.User{}
	o.QueryTable(new(models.User)).Filter("user_name", author).One(&user)
	post := models.Post{
		Title:   title,
		Desc:    desc,
		Content: content,
		Cover:   cover,
		Author:  &user, //只能传结构体对象
	}
	_, err2 := o.Insert(&post)
	if err2 != nil {
		p.Data["json"] = map[string]interface{}{"code": 500, "msg": err2}
		p.ServeJSON()
	}

	p.Data["json"] = map[string]interface{}{"code": 200, "msg": "添加成功"}
	p.ServeJSON()

}

// PostDelete 删除
func (p *PostController) PostDelete() {
	id, err := p.GetInt("id")

	if err != nil {
		p.Ctx.WriteString("id参数错误")
	}

	o := orm.NewOrm()
	_, err2 := o.QueryTable(new(models.Post)).Filter("id", id).Delete()
	if err2 != nil {
		logs.Error(err2)
		p.Ctx.WriteString("删除失败")
	}
	p.Redirect(beego.URLFor("PostController.Get"), 302)
}

// ToEdit 展示修改
func (p *PostController) ToEdit() {
	id, err := p.GetInt("id")
	if err != nil {
		p.Ctx.WriteString("id参数错误")
	}

	o := orm.NewOrm()
	var post models.Post
	o.QueryTable(new(models.Post)).Filter("id", id).One(&post)
	p.Data["post"] = post
	p.TplName = "cms/post-edit.html"
}

// DoEdit 提交修改
func (p *PostController) DoEdit() {
	o := orm.NewOrm()

	id, err := p.GetInt("id")
	if err != nil {
		//p.Ctx.WriteString("id参数错误")，我们使用的post方法，那么就用ajax吧
		p.Data["json"] = map[string]interface{}{"code": 500, "msg": "id参数错误"}
		p.ServeJSON()
	} else {
		qs := o.QueryTable(new(models.Post)).Filter("id", id)
		title := p.GetString("title")
		desc := p.GetString("desc")
		content := p.GetString("content")

		f, h, err1 := p.GetFile("cover")
		defer func(f multipart.File) {
			err := f.Close()
			if err != nil {
				logs.Error(err)
			}
		}(f)
		if err1 != nil {
			//err1报错，说明更改了文件但是没有给它上传文件，这时候我们可以忽略该错误就可以了，直接执行更新就行
			_, err2 := qs.Update(orm.Params{
				"title":   title,
				"desc":    desc,
				"content": content,
			})
			if err2 != nil {
				p.Data["json"] = map[string]interface{}{"code": 500, "msg": "更新失败"}
				p.ServeJSON()
			}
			p.Data["json"] = map[string]interface{}{"code": 200, "msg": "更新成功"}
			p.ServeJSON()

		} else {
			//这时说明用户是上传了新文件，我们需要进行更新
			//格式化名字
			//生成时间戳，防止重名
			timeUnix := time.Now().Unix()               //int64
			time_str := strconv.FormatInt(timeUnix, 10) //int64转换为字符串
			//保存获取到的文件
			path := "static/upload/user/" + time_str + h.Filename
			//把上传的文件保存到本地
			err3 := p.SaveToFile("cover", path)
			if err3 != nil {
				p.Data["json"] = map[string]interface{}{"code": 500, "msg": "保存文件错误"}
				p.ServeJSON()
			} else {
				//更新
				_, err4 := qs.Update(orm.Params{
					"title":   title,
					"desc":    desc,
					"content": content,
					"cover":   path,
				})
				if err4 != nil {
					p.Data["json"] = map[string]interface{}{"code": 500, "msg": "更新失败"}
					p.ServeJSON()
				} else {
					p.Data["json"] = map[string]interface{}{"code": 200, "msg": "更新成功"}
					p.ServeJSON()
				}
			}
		}
	}
}
