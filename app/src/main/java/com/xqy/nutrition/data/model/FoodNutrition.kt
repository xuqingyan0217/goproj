package com.xqy.nutrition.data.model

/**
 * 食物营养信息数据模型
 */
data class FoodNutrition(
    val name: String,        // 食物名称
    val calories: Float,     // 卡路里
    val protein: Float,      // 蛋白质
    val fat: Float,          // 脂肪
    val carbs: Float = 0f    // 碳水化合物，默认为0
)