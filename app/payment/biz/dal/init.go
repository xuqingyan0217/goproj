package dal

import (
	"gomall/app/payment/biz/dal/mysql"
	"gomall/app/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
