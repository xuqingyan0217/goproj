package front

import (
	"Content/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CommentController struct {
	beego.Controller
}

func (c *CommentController) Post() {
	//获取到当前选择文章的id
	post_id, err := c.GetInt("post_id")
	if err != nil {
		fmt.Println(err)
		c.Data["json"] = map[string]interface{}{"code": 500, "msg": "id参数错误"}
		c.ServeJSON()
		return
	} else {
		o := orm.NewOrm()
		post := models.Post{}
		o.QueryTable(new(models.Post)).Filter("id", post_id).One(&post)
		//获取当前用户评论框里面输入的内容
		content := c.GetString("content")
		//获取当前用户名
		user_name := c.GetSession("cms_user_name")
		//通过用户名去查询并赋值给结构体，这是因为Comment模型下需要User对象，而不是一个单纯的值
		user := models.User{}
		o.QueryTable(new(models.User)).Filter("user_name", user_name).One(&user)
		//获取p_id，这个用来判断当前评论是一级还是二级
		pid, err1 := c.GetInt("pid")
		if err1 != nil {
			//一级评论是默认的，所以这个pid是获取不到的,会报错，我们赋值一级id就行了
			pid = 0
		}
		//是二级评论，能直接获取到
		comment := models.Comment{
			Post:    &post,
			Content: content,
			PId:     pid,
			Author:  &user,
		}
		_, err3 := o.Insert(&comment)
		if err3 != nil {
			c.Data["json"] = map[string]interface{}{"code": 500, "msg": "评论出错，请重试"}
			c.ServeJSON()
			return
		} else {
			o := orm.NewOrm()
			newComment := models.Comment{}
			// 根据刚插入评论的id（即comment.Id）去查询对应的评论记录
			err := o.QueryTable(new(models.Comment)).Filter("id", comment.Id).One(&newComment)
			if err != nil {
				// 如果查询出现错误，可以考虑返回一个合适的错误提示给前端，或者进行日志记录等操作
				c.Data["json"] = map[string]interface{}{"code": 500, "msg": "获取新评论数据出错，请重试"}
				c.ServeJSON()
				return
			}
			c.Data["json"] = map[string]interface{}{"code": 200, "msg": "评论成功", "comment": newComment}
			c.ServeJSON()
		}
	}

}
