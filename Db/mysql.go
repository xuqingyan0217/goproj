package Db
//这个包用来我们进行数据库mysql和redis的连接
import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"

	"fmt"
)

//连接数据库
func init() {
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	database := beego.AppConfig.String("database")
	//datasource := "root:123456@tcp(localhost:3306)/content?charset=utf8&loc=Local"
	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local",username,password,host,port,database)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	fmt.Println(datasource)
	orm.RegisterDataBase("default", "mysql", datasource)
	fmt.Println("连接成功")
	//自动显示建表信息
	//-db数据库名字，default默认default
	//-force	drop if exists，默认false
	//-v	打印信息
	orm.RunSyncdb("default", false, true)
}
