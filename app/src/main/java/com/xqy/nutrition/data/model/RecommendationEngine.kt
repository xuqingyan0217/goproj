package com.xqy.nutrition.data.model

import android.content.Context
import android.util.Log
import com.xqy.ui.screens.FoodRecommendation
import kotlin.math.abs

/**
 * 食物推荐引擎
 * 基于用户的历史食品记录和营养目标，生成个性化的食物推荐
 */
class RecommendationEngine private constructor(context: Context) {
    
    private val userProfileManager = UserProfileManager.getInstance(context)
    
    companion object {
        private const val TAG = "RecommendationEngine"
        
        // 食物数据库 - 包含各种食物的营养信息
        private val foodDatabase = listOf(
            // 主食类
            FoodRecommendation("全麦面包", 265, 8, 3, 43, 4.5f),
            FoodRecommendation("糙米饭", 216, 5, 2, 45, 4.2f),
            FoodRecommendation("燕麦片", 389, 13, 7, 66, 4.6f),
            FoodRecommendation("藜麦", 120, 4, 2, 21, 4.3f),
            FoodRecommendation("荞麦面", 231, 8, 1, 46, 4.1f),
            FoodRecommendation("红薯", 86, 2, 0, 20, 4.4f),
            FoodRecommendation("玉米", 96, 3, 1, 21, 4.2f),
            
            // 肉类
            FoodRecommendation("鸡胸肉", 165, 31, 3, 0, 4.8f),
            FoodRecommendation("牛里脊", 143, 26, 4, 0, 4.7f),
            FoodRecommendation("猪里脊", 143, 26, 4, 0, 4.5f),
            FoodRecommendation("羊排", 294, 25, 21, 0, 4.6f),
            FoodRecommendation("鸭胸肉", 201, 23, 11, 0, 4.4f),
            FoodRecommendation("火鸡胸肉", 157, 30, 3, 0, 4.5f),
            FoodRecommendation("兔肉", 173, 33, 4, 0, 4.3f),
            
            // 海鲜类
            FoodRecommendation("三文鱼", 208, 22, 13, 0, 4.7f),
            FoodRecommendation("金枪鱼", 184, 30, 6, 0, 4.6f),
            FoodRecommendation("虾", 99, 24, 1, 1, 4.8f),
            FoodRecommendation("扇贝", 111, 23, 1, 3, 4.7f),
            FoodRecommendation("鳕鱼", 105, 23, 1, 0, 4.5f),
            FoodRecommendation("龙虾", 128, 27, 1, 1, 4.9f),
            FoodRecommendation("牡蛎", 169, 19, 5, 4, 4.4f),
            
            // 蔬菜类
            FoodRecommendation("西兰花", 34, 3, 0, 7, 4.7f),
            FoodRecommendation("菠菜", 23, 3, 0, 4, 4.6f),
            FoodRecommendation("胡萝卜", 41, 1, 0, 10, 4.3f),
            FoodRecommendation("芦笋", 20, 2, 0, 4, 4.5f),
            FoodRecommendation("南瓜", 26, 1, 0, 6, 4.2f),
            FoodRecommendation("花椰菜", 25, 2, 0, 5, 4.4f),
            FoodRecommendation("番茄", 18, 1, 0, 4, 4.3f),
            
            // 水果类
            FoodRecommendation("蓝莓", 57, 1, 0, 14, 4.8f),
            FoodRecommendation("草莓", 32, 1, 0, 8, 4.7f),
            FoodRecommendation("苹果", 52, 0, 0, 14, 4.5f),
            FoodRecommendation("香蕉", 89, 1, 0, 23, 4.4f),
            FoodRecommendation("橙子", 47, 1, 0, 12, 4.6f),
            FoodRecommendation("猕猴桃", 61, 1, 0, 15, 4.7f),
            FoodRecommendation("牛油果", 160, 2, 15, 9, 4.9f),
            
            // 坚果类
            FoodRecommendation("杏仁", 579, 21, 50, 22, 4.8f),
            FoodRecommendation("核桃", 654, 15, 65, 14, 4.7f),
            FoodRecommendation("腰果", 553, 18, 44, 30, 4.6f),
            FoodRecommendation("开心果", 562, 20, 45, 28, 4.7f),
            FoodRecommendation("花生", 567, 26, 49, 16, 4.5f),
            
            // 豆类
            FoodRecommendation("黑豆", 341, 22, 1, 62, 4.4f),
            FoodRecommendation("鹰嘴豆", 364, 19, 6, 61, 4.3f),
            FoodRecommendation("红豆", 337, 24, 1, 60, 4.2f),
            FoodRecommendation("豆腐", 144, 17, 9, 3, 4.5f),
            FoodRecommendation("豆浆", 33, 3, 2, 1, 4.3f),
            
            // 乳制品
            FoodRecommendation("希腊酸奶", 59, 10, 0, 3, 4.7f),
            FoodRecommendation("奶酪", 402, 25, 33, 1, 4.6f),
            FoodRecommendation("牛奶", 42, 3, 1, 5, 4.4f),
            FoodRecommendation("酸奶", 59, 3, 3, 4, 4.5f),
            
            // 零食类
            FoodRecommendation("黑巧克力", 598, 8, 43, 46, 4.8f),
            FoodRecommendation("能量棒", 458, 10, 22, 58, 4.3f),
            FoodRecommendation("蛋白棒", 360, 20, 12, 40, 4.5f),
            FoodRecommendation("水果干", 359, 4, 1, 88, 4.2f)
        )
        
        @Volatile
        private var instance: RecommendationEngine? = null
        
        fun getInstance(context: Context): RecommendationEngine {
            return instance ?: synchronized(this) {
                instance ?: RecommendationEngine(context.applicationContext).also { instance = it }
            }
        }
    }
    
    /**
     * 获取个性化食物推荐
     * 基于用户的营养目标和当天已摄入的营养，推荐最适合的食物
     * @return 推荐食物列表
     */
    fun getRecommendations(): List<FoodRecommendation> {
        // 获取用户的营养目标
        val caloriesGoal = userProfileManager.getCalories().toIntOrNull() ?: 2500
        val proteinGoal = userProfileManager.getProtein().toIntOrNull() ?: 90
        val carbsGoal = userProfileManager.getCarbs().toIntOrNull() ?: 225
        val fatGoal = userProfileManager.getFat().toIntOrNull() ?: 60
        
        // 获取用户当天已摄入的营养
        val todayNutrition = FoodRecordManager.getTodayTotalNutrition()
        val consumedCalories = todayNutrition.first
        val consumedProtein = todayNutrition.second.replace("g", "").toFloatOrNull() ?: 0f
        val consumedFat = todayNutrition.third.replace("g", "").toFloatOrNull() ?: 0f
        val consumedCarbs = todayNutrition.fourth.replace("g", "").toFloatOrNull() ?: 0f
        
        // 计算剩余需要摄入的营养
        val remainingCalories = caloriesGoal - consumedCalories
        val remainingProtein = proteinGoal - consumedProtein.toInt()
        val remainingCarbs = carbsGoal - consumedCarbs.toInt()
        val remainingFat = fatGoal - consumedFat.toInt()
        
        Log.d(TAG, "营养目标 - 热量: $caloriesGoal, 蛋白质: $proteinGoal, 碳水: $carbsGoal, 脂肪: $fatGoal")
        Log.d(TAG, "已摄入 - 热量: $consumedCalories, 蛋白质: $consumedProtein, 碳水: $consumedCarbs, 脂肪: $consumedFat")
        Log.d(TAG, "剩余 - 热量: $remainingCalories, 蛋白质: $remainingProtein, 碳水: $remainingCarbs, 脂肪: $remainingFat")
        
        // 根据剩余营养需求对食物进行评分
        val scoredFoods = foodDatabase.map { food ->
            // 计算食物对剩余营养的贡献度
            val caloriesScore = if (remainingCalories > 0) {
                // 如果还需要摄入热量，那么食物的热量越接近剩余热量的一定比例（如20%）越好
                val idealCalories = remainingCalories * 0.2
                1.0 - abs(food.calories - idealCalories) / idealCalories
            } else {
                // 如果已经超过热量目标，那么低热量食物更好
                1.0 - (food.calories / 500.0) // 假设500卡是一餐的标准热量
            }
            
            // 蛋白质评分 - 如果需要蛋白质，高蛋白食物得分高
            val proteinScore = if (remainingProtein > 0) {
                food.protein.toDouble() / 30.0 // 假设30g是一餐的标准蛋白质量
            } else {
                1.0 - (food.protein.toDouble() / 30.0)
            }
            
            // 碳水评分
            val carbsScore = if (remainingCarbs > 0) {
                food.carbs.toDouble() / 50.0 // 假设50g是一餐的标准碳水量
            } else {
                1.0 - (food.carbs.toDouble() / 50.0)
            }
            
            // 脂肪评分
            val fatScore = if (remainingFat > 0) {
                food.fat.toDouble() / 20.0 // 假设20g是一餐的标准脂肪量
            } else {
                1.0 - (food.fat.toDouble() / 20.0)
            }
            
            // 综合评分 - 根据用户的营养需求调整各项的权重
            val totalScore = when {
                remainingProtein > 20 -> {
                    // 如果蛋白质缺口大，增加蛋白质权重
                    caloriesScore * 0.2 + proteinScore * 0.5 + carbsScore * 0.15 + fatScore * 0.15
                }
                remainingCarbs > 50 -> {
                    // 如果碳水缺口大，增加碳水权重
                    caloriesScore * 0.2 + proteinScore * 0.15 + carbsScore * 0.5 + fatScore * 0.15
                }
                remainingFat > 15 -> {
                    // 如果脂肪缺口大，增加脂肪权重
                    caloriesScore * 0.2 + proteinScore * 0.15 + carbsScore * 0.15 + fatScore * 0.5
                }
                else -> {
                    // 平衡配置
                    caloriesScore * 0.25 + proteinScore * 0.25 + carbsScore * 0.25 + fatScore * 0.25
                }
            }
            
            // 考虑食物原始评分
            val finalScore = totalScore * 0.8 + (food.rating / 5.0) * 0.2
            
            Pair(food, finalScore)
        }
        
        // 根据评分排序并返回前15个推荐
        return scoredFoods.sortedByDescending { it.second }.take(15).map { it.first }
    }

    /**
     * 分析用户的饮食习惯
     * @return 饮食习惯分析结果
     */
    fun analyzeEatingHabits(): String {
        val allRecords = FoodRecordManager.allRecords
        if (allRecords.isEmpty()) {
            return "暂无足够的饮食记录进行分析"
        }
        
        // 计算平均每日营养摄入
        val totalCalories = allRecords.sumOf { it.calories.toDouble() }
        val totalProtein = allRecords.sumOf { it.protein.toDouble() }
        val totalFat = allRecords.sumOf { it.fat.toDouble() }
        val totalCarbs = allRecords.sumOf { it.carbs.toDouble() }
        
        // 获取不同日期的数量
        val uniqueDates = allRecords.map { it.date }.toSet().size
        
        // 如果没有足够的数据，返回默认消息
        if (uniqueDates < 3) {
            return "需要至少3天的饮食记录才能进行分析"
        }
        
        val avgCalories = totalCalories / uniqueDates
        val avgProtein = totalProtein / uniqueDates
        val avgFat = totalFat / uniqueDates
        val avgCarbs = totalCarbs / uniqueDates
        
        // 获取用户的营养目标
        val caloriesGoal = userProfileManager.getCalories().toIntOrNull() ?: 2500
        val proteinGoal = userProfileManager.getProtein().toIntOrNull() ?: 90
        val carbsGoal = userProfileManager.getCarbs().toIntOrNull() ?: 225
        val fatGoal = userProfileManager.getFat().toIntOrNull() ?: 60
        
        // 分析结果
        val caloriesAnalysis = when {
            avgCalories > caloriesGoal * 1.1 -> "热量摄入过高"
            avgCalories < caloriesGoal * 0.9 -> "热量摄入不足"
            else -> "热量摄入适中"
        }
        
        val proteinAnalysis = when {
            avgProtein > proteinGoal * 1.1 -> "蛋白质摄入过高"
            avgProtein < proteinGoal * 0.9 -> "蛋白质摄入不足"
            else -> "蛋白质摄入适中"
        }
        
        val carbsAnalysis = when {
            avgCarbs > carbsGoal * 1.1 -> "碳水摄入过高"
            avgCarbs < carbsGoal * 0.9 -> "碳水摄入不足"
            else -> "碳水摄入适中"
        }
        
        val fatAnalysis = when {
            avgFat > fatGoal * 1.1 -> "脂肪摄入过高"
            avgFat < fatGoal * 0.9 -> "脂肪摄入不足"
            else -> "脂肪摄入适中"
        }
        
        return "饮食习惯分析:\n" +
                "- $caloriesAnalysis (平均 ${avgCalories.toInt()} 千卡/天)\n" +
                "- $proteinAnalysis (平均 ${avgProtein.toInt()} 克/天)\n" +
                "- $carbsAnalysis (平均 ${avgCarbs.toInt()} 克/天)\n" +
                "- $fatAnalysis (平均 ${avgFat.toInt()} 克/天)"
    }
}