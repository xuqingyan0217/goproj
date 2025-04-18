package model

import (
	"context"
	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`

	Products []Product `json:"products" gorm:"many2many:product_category;"`
}

func (Category) TableName() string {
	return "category"
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (c CategoryQuery) GetProductsByCategoryName(name string) (categories []Category, err error) {
	// preload预加载一下products,结果会显示到前端
	err = c.db.WithContext(c.ctx).Model(&Category{}).
		Where(&Category{Name: name}).Preload("Products").
		Find(&categories).Error
	return categories, err
}

func (c CategoryQuery) GetCategoriesByProductId(productId int) (categoryNames []string, err error) {
	var categories []Category
	// 使用GORM的关联查询特性
	err = c.db.WithContext(c.ctx).Model(&Category{}).
		Preload("Products", "id = ?", productId).
		Where("id IN (SELECT category_id FROM product_category WHERE product_id = ?)", productId).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 提取类别名称
	categoryNames = make([]string, len(categories))
	for i, category := range categories {
		categoryNames[i] = category.Name
	}
	return categoryNames, nil
}

func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{
		ctx: ctx,
		db:  db,
	}
}

func (c CategoryQuery) GetAllCategoryNames() (categoryNames []string, err error) {
	var categories []Category
	// 查询所有类别
	err = c.db.WithContext(c.ctx).Model(&Category{}).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 提取类别名称
	categoryNames = make([]string, len(categories))
	for i, category := range categories {
		categoryNames[i] = category.Name
	}
	return categoryNames, nil
}
