package models

import "time"

type CommentTree struct {
	Id 	int
	Content string
	Author  *User
	CreateTime	time.Time
	Children  []*CommentTree
}

type Comment struct {
	Id 		int `orm:"pk;auto"`
	CreateTime	time.Time	`orm:"auto_now_add;type(datetime);description(创建时间)"`
	Content		string  `orm:"size(4000);description(评论的内容)"`
	//PId是我们判断当前评论是属于一级评论还是二级评论的标识，我们默认是0（一级评论）
	PId 		int		`orm:"description(父级评论);default(0)"`
	//建立起和用户表的关系，一个用户可以发表多条评论
	Author 	*User `orm:"rel(fk);description(评论作者)"`
	Post	*Post	`orm:"rel(fk);description(帖子外键)"`
}

//评论的层级关系，下面的意思是0代表是一级评论，而第二条的pid是1，说明是第一条的二级评论
//1 0
//2 1

func (c *Comment) TableName() string {
	return "sys_post_comment"
}
