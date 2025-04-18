package dal

import (
	"gomall/app/order/biz/dal/mysql"
	"gomall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
