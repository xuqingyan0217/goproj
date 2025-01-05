package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id 			int `orm:"pk;auto"`
	UserName 	string	`orm:"description(用户名);index;unique"`
	Password 	string	`orm:"description(密码)"`
	IsAdmin 	int		`orm:"description(1是管理员，2是普通用户);default(2)"`
	CreateTime	time.Time	`orm:"auto_now_add;type(datetime);description(创建时间)"`
	Cover 		string	`orm:"description(头像);default(static/upload/bq3.png)"`
	//建立和帖子模型之间的关系
	Posts		[]*Post `orm:"description(用户的帖子切片);reverse(many)"`
	//建立和评论comment模型之间的关系
	Comments	[]*Comment `orm:"description(用户的评论);reverse(many)"`
}

//TableName 修改结构体在模型中的表名
func (u *User) TableName() string  {
	return "sys_user"
}

func init()  {
	orm.RegisterModel(new(User), new(Post), new(Comment))
}