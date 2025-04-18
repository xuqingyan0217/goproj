package com.xqy.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Star
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.xqy.nutrition.data.model.FoodRecordManager
import com.xqy.nutrition.data.model.RecommendationEngine
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.background

data class FoodRecommendation(
    val name: String,
    val calories: Int,
    val protein: Int,
    val fat: Int,
    val carbs: Int,
    val rating: Float
)

@Composable
fun RecommendScreen() {
    val context = LocalContext.current
    val recommendationEngine = remember { RecommendationEngine.getInstance(context) }
    
    // 使用FoodRecordManager的记录变化来触发推荐刷新
    val recommendations = remember(FoodRecordManager.records) {
        recommendationEngine.getRecommendations()
    }
    
    // 获取饮食习惯分析
    val eatingHabitsAnalysis = remember(FoodRecordManager.allRecords) {
        recommendationEngine.analyzeEatingHabits()
    }
    
    // 获取当前营养摄入情况
    val todayNutrition = remember(FoodRecordManager.records) {
        FoodRecordManager.getTodayTotalNutrition()
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        item {
            Text(
                text = "智能推荐",
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold
            )
            
            Spacer(modifier = Modifier.height(8.dp))
            
            Text(
                text = "基于您的营养需求和历史记录，为您推荐以下食物",
                fontSize = 14.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            
            // 显示当前营养摄入情况
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 8.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.7f)
                )
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "今日营养摄入",
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Bold
                    )
                    
                    Spacer(modifier = Modifier.height(8.dp))
                    
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceEvenly
                    ) {
                        NutrientSummary("热量", "${todayNutrition.first}千卡", MaterialTheme.colorScheme.primary)
                        NutrientSummary("蛋白质", todayNutrition.second, MaterialTheme.colorScheme.secondary)
                        NutrientSummary("脂肪", todayNutrition.third, MaterialTheme.colorScheme.tertiary)
                        NutrientSummary("碳水", todayNutrition.fourth, MaterialTheme.colorScheme.error)
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // 添加饮食习惯分析卡片
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 8.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.secondaryContainer.copy(alpha = 0.7f)
                )
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "饮食习惯分析",
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Bold
                    )
                    
                    Spacer(modifier = Modifier.height(8.dp))
                    
                    if (eatingHabitsAnalysis.startsWith("需要至少") || eatingHabitsAnalysis.startsWith("暂无足够")) {
                        // 显示提示信息
                        Text(
                            text = eatingHabitsAnalysis,
                            fontSize = 14.sp,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    } else {
                        // 将分析结果拆分为多行并显示
                        val analysisLines = eatingHabitsAnalysis.split("\n")
                        analysisLines.drop(1).forEach { line ->
                            Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.padding(vertical = 4.dp)) {
                                val color = when {
                                    line.contains("过高") -> MaterialTheme.colorScheme.error
                                    line.contains("不足") -> MaterialTheme.colorScheme.tertiary
                                    else -> MaterialTheme.colorScheme.primary
                                }
                                
                                // 添加小圆点
                                Box(
                                    modifier = Modifier
                                        .size(8.dp)
                                        .background(color = color, shape = CircleShape)
                                )
                                
                                Spacer(modifier = Modifier.width(8.dp))
                                
                                Text(
                                    text = line.trim().removePrefix("-").trim(),
                                    fontSize = 14.sp,
                                    color = MaterialTheme.colorScheme.onSurface
                                )
                            }
                        }
                        
                        // 添加改进建议
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            text = "改进建议",
                            fontSize = 14.sp,
                            fontWeight = FontWeight.Medium,
                            color = MaterialTheme.colorScheme.onSurface
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        
                        val suggestions = mutableListOf<String>()
                        if (analysisLines.any { it.contains("热量摄入过高") }) {
                            suggestions.add("减少高热量食物摄入，增加体育活动")
                        }
                        if (analysisLines.any { it.contains("热量摄入不足") }) {
                            suggestions.add("适当增加进食量，选择营养密度高的食物")
                        }
                        if (analysisLines.any { it.contains("蛋白质摄入不足") }) {
                            suggestions.add("增加瘦肉、鱼、蛋、豆制品等优质蛋白来源")
                        }
                        if (analysisLines.any { it.contains("脂肪摄入过高") }) {
                            suggestions.add("减少油炸食品和高脂肪食物摄入")
                        }
                        if (analysisLines.any { it.contains("碳水摄入过高") }) {
                            suggestions.add("减少精制碳水化合物，选择全谷物食品")
                        }
                        
                        suggestions.forEach { suggestion ->
                            Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.padding(vertical = 2.dp)) {
                                Box(
                                    modifier = Modifier
                                        .size(6.dp)
                                        .background(color = MaterialTheme.colorScheme.secondary, shape = CircleShape)
                                )
                                Spacer(modifier = Modifier.width(8.dp))
                                Text(
                                    text = suggestion,
                                    fontSize = 13.sp,
                                    color = MaterialTheme.colorScheme.onSurfaceVariant
                                )
                            }
                        }
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
        }
        
        items(recommendations) { recommendation ->
            RecommendationCard(recommendation)
        }
    }
    }


@Composable
private fun NutrientSummary(label: String, value: String, color: Color) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(
            text = value,
            fontSize = 16.sp,
            fontWeight = FontWeight.Bold,
            color = color
        )
        Text(
            text = label,
            fontSize = 14.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )
    }
}

@Composable
private fun RecommendationCard(recommendation: FoodRecommendation) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = recommendation.name,
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold
                )
                
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(
                        imageVector = Icons.Default.Star,
                        contentDescription = "评分",
                        tint = Color.Red,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = recommendation.rating.toString(),
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.primary
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                NutrientInfo(
                    value = recommendation.calories,
                    unit = "千卡",
                    label = "热量",
                    color = MaterialTheme.colorScheme.primary
                )
                NutrientInfo(
                    value = recommendation.protein,
                    unit = "g",
                    label = "蛋白质",
                    color = MaterialTheme.colorScheme.secondary
                )
                NutrientInfo(
                    value = recommendation.fat,
                    unit = "g",
                    label = "脂肪",
                    color = MaterialTheme.colorScheme.tertiary
                )
                NutrientInfo(
                    value = recommendation.carbs,
                    unit = "g",
                    label = "碳水",
                    color = MaterialTheme.colorScheme.error
                )
            }
        }
    }
}

@Composable
private fun NutrientInfo(
    value: Int,
    unit: String,
    label: String,
    color: Color
) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = value.toString(),
            fontSize = 24.sp,
            fontWeight = FontWeight.Bold,
            color = color
        )
        Text(
            text = unit,
            fontSize = 14.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )
        Text(
            text = label,
            fontSize = 16.sp,
            fontWeight = FontWeight.Medium
        )
    }
}