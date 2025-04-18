package model

import (
	"context"
	"gorm.io/gorm"
)

type Consignee struct {
	Email         string
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}
type Order struct {
	gorm.Model
	OrderId      string `gorm:"type:varchar(255);uniqueIndex"`
	UserId       uint32 `gorm:"type:int(11)"`
	UserCurrency string `gorm:"type:varchar(10)"`
	// embedded表示嵌入一个结构体
	Consignee Consignee `gorm:"embedded"`
	// 关联一对多，一个订单id下可能会有多个订单项
	OrderItems []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
}

func (Order) TableName() string {
	return "order"
}

func ListOrder(ctx context.Context, db *gorm.DB, userId uint32) ([]*Order, error) {
	var orders []*Order
	err := db.WithContext(ctx).Where("user_id = ?", userId).Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}
