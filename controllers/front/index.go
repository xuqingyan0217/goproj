package front

import (
	"Content/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type IndexController struct {
	beego.Controller
}

func (i *IndexController) Get() {
	o := orm.NewOrm()

	// 获取请求参数
	page, _ := i.GetInt("page", 1)    // 页码，默认第1页
	keyword := i.GetString("keyword") // 搜索关键词
	pageSize := 5                     // 每页显示数量

	// 构建基础查询
	query := o.QueryTable(new(models.Post)).RelatedSel()

	// 如果有关键词，添加模糊搜索
	if keyword != "" {
		query = query.Filter("Title__icontains", keyword)
	}

	// 查询总记录数
	totalCount, _ := query.Count()

	// 分页查询
	var posts []models.Post
	_, err := query.
		Limit(pageSize, (page-1)*pageSize).
		OrderBy("-CreateTime"). // 按创建时间倒序
		All(&posts)

	if err != nil {
		// 处理查询错误
		i.Data["error"] = "查询失败"
	}
	// 判断是否还要显示"显示更多"
	var hasMore = false
	if (int64(page)-1)*int64(pageSize)+int64(len(posts)) < totalCount {
		hasMore = true
	}
	// 区分是普通页面加载还是AJAX加载
	isAjax := i.GetString("_ajax") == "true"

	// 设置模板数据
	i.Data["posts"] = posts
	i.Data["keyword"] = keyword
	i.Data["currentPage"] = page
	i.Data["hasMore"] = hasMore

	// 根据请求类型返回不同响应
	if isAjax {
		// 如果是AJAX请求，只返回文章列表HTML片段
		i.TplName = "front/posts_list.html"
	} else {
		// 普通页面加载，返回完整页面
		i.TplName = "front/index.html"
	}
}

func (i *IndexController) PostDetail() {
	id, _ := i.GetInt("id")
	o := orm.NewOrm()
	post := models.Post{}
	o.QueryTable(new(models.Post)).RelatedSel().Filter("id", id).One(&post)

	//评论模块，依据当前详情页的id去查找对应的评论
	comments := []models.Comment{}
	//这里不但要过滤只是该帖子id下的评论，还要过滤到除一级评论（0）以外的评论，保证我们最终结果只针对一层评论
	o.QueryTable(new(models.Comment)).RelatedSel().Filter("post_id", id).Filter("p_id", 0).All(&comments)
	//定义结构体切片，用于把这个树连起来
	comment_trees := []models.CommentTree{}

	for _, comment := range comments {
		pid := comment.Id
		comment_tree := models.CommentTree{
			Id:         comment.Id,
			Content:    comment.Content,
			Author:     comment.Author,
			CreateTime: comment.CreateTime,
			Children:   []*models.CommentTree{},
		}
		//递归
		GetChild(pid, &comment_tree)
		//连接起来
		comment_trees = append(comment_trees, comment_tree)
	}
	i.Data["comment_trees"] = comment_trees
	i.Data["post"] = post
	i.TplName = "front/detail.html"
}

// GetChild 递归函数
func GetChild(pid int, comment_tree *models.CommentTree) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Comment))
	comments := []models.Comment{}
	_, err := qs.Filter("p_id", pid).RelatedSel().All(&comments)
	if err != nil {
		return
	}

	//查询二级以及更高层级的评论
	for i := 0; i < len(comments); i++ {
		//上面的pid是一级评论下的pid，而后续的递归需要获取二级评论下的pid，所以这里需要更新pid，这样往下递归时保证无差错
		pid := comments[i].Id
		child := models.CommentTree{
			Id:         comments[i].Id,
			Content:    comments[i].Content,
			Author:     comments[i].Author,
			CreateTime: comments[i].CreateTime,
			Children:   []*models.CommentTree{},
		}
		comment_tree.Children = append(comment_tree.Children, &child)
		GetChild(pid, &child)
	}
	return
}

//数据结构
//[
//	{
//		id : 1,
//		content : "xxx",
//		children : [
//			{
//				id : 2,
//				content : "xxx",
//				children : [
//					{
//					id : 3,
//					content : "xxx",
//					}
//				]
//			}
//		]
//	}
//]
