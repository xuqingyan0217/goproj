package dal

import (
	"gomall/app/cart/biz/dal/mysql"
	"gomall/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
