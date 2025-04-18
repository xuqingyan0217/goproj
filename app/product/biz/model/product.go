package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Picture     string         `json:"picture"`
	Price       float32        `json:"price"`
	DeleteAt    gorm.DeletedAt `gorm:"index"`
	Categories  []Category     `json:"categories" gorm:"many2many:product_category;"`
}

func (Product) TableName() string {
	return "product"
}

// ProductQuery ProductMutation读写分离
type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}
type ProductMutation struct {
	ctx         context.Context
	db          *gorm.DB
	cacheClient *redis.Client
}

func (p ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Unscoped().Model(&Product{}).First(&product, productId).Error
	return product, err
}

func (p ProductQuery) SearchProducts(q string) (products []*Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Find(&products,
		"name like ? or description like ?", "%"+q+"%", "%"+q+"%").Error
	return products, err
}

func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{ctx: ctx, db: db}
}

func NewProductMutation(ctx context.Context, db *gorm.DB, cacheClient *redis.Client) *ProductMutation {
	return &ProductMutation{
		ctx:         ctx,
		db:          db,
		cacheClient: cacheClient,
	}
}

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
}

// GetById 根据产品ID获取产品信息。
// 首先尝试从缓存中获取产品信息，如果缓存中不存在，则从数据库中获取，并将结果缓存。
// 参数:
//
//	productId - 产品的ID。
//
// 返回值:
//
//	product - 产品信息。
//	err - 错误信息，如果获取产品信息时发生错误，则返回该错误。
func (c CachedProductQuery) GetById(productId int) (product Product, err error) {
	// 构造缓存键。
	cacheKey := fmt.Sprintf("%s_%s_%d", c.prefix, "product_by_id", productId)

	// 从缓存中获取产品信息。
	cacheResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	// 解析缓存结果。
	err = func() error {
		// 检查缓存获取操作是否有错误。
		if err := cacheResult.Err(); err != nil {
			return err
		}

		// 将缓存结果转换为字节切片。
		cacheResultByte, err := cacheResult.Bytes()
		if err != nil {
			return err
		}

		// 解析JSON格式的缓存数据。
		err = json.Unmarshal(cacheResultByte, &product)
		if err != nil {
			return err
		}

		return nil
	}()
	// 如果缓存中没有数据或者解析缓存数据时发生错误，则从数据库中获取产品信息。
	if err != nil {
		// 从数据库中获取产品信息。
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return Product{}, err
		}

		// 将获取的产品信息编码为JSON格式。
		encoded, err := json.Marshal(product)
		if err != nil {
			return Product{}, err
		}

		// 将从数据库中获取的产品信息缓存起来，设置缓存过期时间为1小时。
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}

	// 返回获取的产品信息。
	return
}

func (c CachedProductQuery) SearchProducts(q string) (products []*Product, err error) {
	// 假设命中率很低，直接调用数据库的查询方法。
	return c.productQuery.SearchProducts(q)
}

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, cacheClient *redis.Client) *CachedProductQuery {
	return &CachedProductQuery{
		productQuery: *NewProductQuery(ctx, db),
		cacheClient:  cacheClient,
		prefix:       "shop",
	}
}

// CreateProduct 创建产品。
// 参数:
//
//	product - 产品信息。
//
// 返回值:
//
//	product - 创建的产品信息。
//	err - 错误信息，如果创建产品信息时发生错误，则返回该错误。
func (m ProductMutation) CreateProduct(input Product) (product Product, err error) {
	// 开启事务
	err = m.db.WithContext(m.ctx).Transaction(func(tx *gorm.DB) error {
		// 保存产品的类别信息
		categories := input.Categories

		// 创建产品基本信息
		input.Categories = nil
		if err := tx.Create(&input).Error; err != nil {
			return err
		}

		// 处理类别关联
		if len(categories) > 0 {
			// 查找已存在的类别
			var existingCategories []Category
			categoryNames := extractCategoryNames(categories)
			if err := tx.Where("name IN ?", categoryNames).Find(&existingCategories).Error; err != nil {
				return err
			}

			// 使用 map 存储已存在类别的名称
			existingCategoryMap := make(map[string]Category)
			for _, existing := range existingCategories {
				existingCategoryMap[existing.Name] = existing
			}

			// 存储新的类别
			for _, category := range categories {
				if _, exists := existingCategoryMap[category.Name]; !exists {
					// 创建新类别
					newCategory := Category{
						Name:        category.Name,
						Description: category.Description,
					}
					if err := tx.Create(&newCategory).Error; err != nil {
						return err
					}
					existingCategoryMap[newCategory.Name] = newCategory
				}
			}

			// 建立商品和类别的关联关系
			for _, category := range categories {
				if existingCategory, ok := existingCategoryMap[category.Name]; ok {
					if err := tx.Model(&input).Association("Categories").Append(&existingCategory); err != nil {
						return err
					}
				}
			}
		}

		// 重新加载商品信息，包括关联的类别
		if err := tx.Preload("Categories").First(&product, input.ID).Error; err != nil {
			return err
		}

		return nil
	})

	return product, err
}

// 辅助函数：提取类别名称
func extractCategoryNames(categories []Category) []string {
	var categoryNames []string
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Name)
	}
	return categoryNames
}

// UpdateProduct 更新产品。
// 参数:
//
//	productId - 产品的ID。
//	input - 产品信息。
//
// 返回值:
//
//	err - 错误信息，如果更新产品信息时发生错误，则返回该错误。
func (m ProductMutation) UpdateProduct(productId int, input Product) (err error) {
	// 开启事务
	return m.db.WithContext(m.ctx).Transaction(func(tx *gorm.DB) error {
		// 更新产品基本信息
		if err := tx.Model(&Product{}).Where("id = ?", productId).Updates(map[string]interface{}{
			"name":        input.Name,
			"description": input.Description,
			"picture":     input.Picture,
			"price":       input.Price,
		}).Error; err != nil {
			return err
		}

		// 如果提供了新的类别信息，更新类别关联
		if len(input.Categories) > 0 {
			// 先清除现有的所有类别关联
			if err := tx.Model(&Product{Base: Base{ID: productId}}).Association("Categories").Clear(); err != nil {
				return err
			}

			// 查找已存在的类别
			var existingCategories []Category
			if err := tx.Where("name IN ?", extractCategoryNames(input.Categories)).Find(&existingCategories).Error; err != nil {
				return err
			}

			// 不存在的类别加入
			for _, category := range input.Categories {
				found := false
				for _, existing := range existingCategories {
					if existing.Name == category.Name {
						found = true
						fmt.Println("existing", existing)
					}
				}
				if !found {
					existingCategories = append(existingCategories, category)
				}
			}

			// 建立新的类别关联
			if err := tx.Model(&Product{Base: Base{ID: productId}}).Association("Categories").Append(existingCategories); err != nil {
				return err
			}
		}

		// 检查类别关联
		if err != nil {
			return err
		}

		// 删除缓存，确保下次获取时能拿到最新数据
		cacheKey := fmt.Sprintf("%s_%s_%d", "shop", "product_by_id", productId)
		if m.cacheClient != nil {
			_ = m.cacheClient.Del(m.ctx, cacheKey)
		}

		return nil
	})
}

// DeleteProduct 删除产品。
// 参数:
//
//	productId - 产品的ID。
//
// 返回值:
//
//	err - 错误信息，如果删除产品信息时发生错误，则返回该错误。
func (m ProductMutation) DeleteProduct(productId int) (err error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("%s_%s_%d", "shop", "product_by_id", productId)

	// 开启事务
	return m.db.WithContext(m.ctx).Transaction(func(tx *gorm.DB) error {
		// 先清除商品与类别的关联关系
		if err := tx.Model(&Product{Base: Base{ID: productId}}).Association("Categories").Clear(); err != nil {
			return err
		}

		// 使用软删除机制删除商品
		if err := tx.Model(&Product{}).Where("id = ?", productId).Delete(&Product{}).Error; err != nil {
			return err
		}

		// 删除缓存，确保下次获取时能拿到最新数据
		if m.cacheClient != nil {
			_ = m.cacheClient.Del(m.ctx, cacheKey)
		}

		return nil
	})
}
