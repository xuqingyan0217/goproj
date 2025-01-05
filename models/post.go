package models

import (
	"time"
)

type Tag struct {
	Id    int
	TName string
}

type Post struct {
	Id      int    `orm:"pk;auto"`
	Title   string `orm:"description(帖子标题)"`
	Desc    string `orm:"description(帖子描述)"`
	Content string `orm:"size(4000);description(帖子内容)"`
	//封面的默认值是暂无图片
	Cover      string     `orm:"description(帖子封面图);default(static/upload/no_pic.jpg)"`
	ReadNum    int        `orm:"description(阅读量);default(0)"`
	StarNum    int        `orm:"description(帖子点赞数);default(0)"`
	CreateTime time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)"`
	Author     *User      `orm:"description(帖子作者);rel(fk)"`
	Comments   []*Comment `orm:"reverse(many)"`
}

// TableName 给模型修改其在结构体里面的表名
func (p *Post) TableName() string {
	return "sys_post"
}

//注册模型，写到user里面就可以了，毕竟都在一个包下面，不必多执行一次init
