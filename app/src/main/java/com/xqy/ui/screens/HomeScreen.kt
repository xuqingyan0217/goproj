/**
 * 主页面
 * 显示用户的每日营养摄入情况，包括卡路里、蛋白质和脂肪的统计数据
 * 提供日期选择功能和最近饮食记录列表
 */
package com.xqy.ui.screens

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material3.*
import androidx.compose.material3.DatePicker
import androidx.compose.material3.DatePickerDialog
import androidx.compose.material3.rememberDatePickerState
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.xqy.nutrition.data.model.UserProfileManager
import com.xqy.ui.components.RecentRecordsList
import com.xqy.ui.components.NutritionItem
import java.time.ZoneId
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import java.util.Locale

/**
 * 主页面组件
 * 包含以下功能：
 * - 显示页面标题
 * - 日期选择器
 * - 营养摄入数据卡片
 * - 卡路里目标进度条
 * - 最近饮食记录列表
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen() {
    val context = LocalContext.current
    val userProfileManager = remember { UserProfileManager.getInstance(context) }
    
    // 添加滚动状态
    val scrollState = rememberScrollState()
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(scrollState) // 添加垂直滚动
            .padding(16.dp)
    ) {
        // 页面标题
        Text(
            text = "营养健康小助手",
            fontSize = 24.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.onSurface,
            modifier = Modifier.padding(bottom = 24.dp)
        )

        // 营养跟踪标题
        Text(
            text = "营养跟踪",
            fontSize = 18.sp,
            fontWeight = FontWeight.SemiBold,
            color = MaterialTheme.colorScheme.onSurface,
            modifier = Modifier.padding(bottom = 12.dp)
        )

        // 日期选择状态管理
        val selectedDate = remember { mutableStateOf(com.xqy.nutrition.data.model.FoodRecordManager.getSelectedDate()) }
        val showDatePicker = remember { mutableStateOf(false) }
        
        // 日期选择器触发按钮
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier
                .clickable { showDatePicker.value = true }
                .padding(bottom = 20.dp)
        ) {
            Text(
                text = com.xqy.nutrition.data.model.FoodRecordManager.formatDateForDisplay(selectedDate.value),
                fontSize = 14.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(modifier = Modifier.width(4.dp))
            Icon(
                imageVector = androidx.compose.material.icons.Icons.Filled.DateRange,
                contentDescription = "选择日期",
                tint = MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.size(16.dp)
            )
        }
        
        // 日期选择器对话框
        if (showDatePicker.value) {
            val datePickerState = rememberDatePickerState()
            DatePickerDialog(
                onDismissRequest = { showDatePicker.value = false },
                confirmButton = {
                    TextButton(
                        onClick = {
                            datePickerState.selectedDateMillis?.let { millis ->
                                val localDate = java.time.Instant.ofEpochMilli(millis)
                                    .atZone(ZoneId.of("Asia/Shanghai"))
                                    .toLocalDate()
                                selectedDate.value = localDate
                                com.xqy.nutrition.data.model.FoodRecordManager.filterRecordsByDate(localDate)
                            }
                            showDatePicker.value = false
                        }
                    ) {
                        Text("确定")
                    }
                },
                dismissButton = {
                    TextButton(onClick = { showDatePicker.value = false }) {
                        Text("取消")
                    }
                }
            ) {
                DatePicker(state = datePickerState)
            }
        }

        // 获取当日营养数据
        val todayNutrition = remember(com.xqy.nutrition.data.model.FoodRecordManager.records, selectedDate.value) {
            com.xqy.nutrition.data.model.FoodRecordManager.getTodayTotalNutrition()
        }

        // 营养数据卡片
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 8.dp),
            colors = CardDefaults.cardColors(
                containerColor = MaterialTheme.colorScheme.surfaceVariant
            )
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                NutritionItem(todayNutrition.first.toString(), "卡路里", MaterialTheme.colorScheme.primary)
                NutritionItem(todayNutrition.second, "蛋白质", MaterialTheme.colorScheme.secondary)
                NutritionItem(todayNutrition.third, "脂肪", MaterialTheme.colorScheme.tertiary)
                NutritionItem(todayNutrition.fourth, "碳水", MaterialTheme.colorScheme.error)
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // 获取用户设置的营养目标
        val targetCalories = userProfileManager.getCalories().toFloat()
        val targetProtein = userProfileManager.getProtein().toFloat()
        val targetFat = userProfileManager.getFat().toFloat()
        val targetCarbs = userProfileManager.getCarbs().toFloat()
        
        // 卡路里目标进度
        val todayCalories = todayNutrition.first
        val caloriesProgress = (todayCalories / targetCalories).coerceIn(0f, 1f)
        val caloriesPercentage = (caloriesProgress * 100).toInt()
        
        // 蛋白质目标进度
        val todayProtein = todayNutrition.second.replace("g", "").toFloat()
        val proteinProgress = (todayProtein / targetProtein).coerceIn(0f, 1f)
        val proteinPercentage = (proteinProgress * 100).toInt()
        
        // 脂肪目标进度
        val todayFat = todayNutrition.third.replace("g", "").toFloat()
        val fatProgress = (todayFat / targetFat).coerceIn(0f, 1f)
        val fatPercentage = (fatProgress * 100).toInt()
        
        // 碳水目标进度
        val todayCarbs = todayNutrition.fourth.replace("g", "").toFloat()
        val carbsProgress = (todayCarbs / targetCarbs).coerceIn(0f, 1f)
        val carbsPercentage = (carbsProgress * 100).toInt()
        
        // 营养目标进度标题
        Text(
            text = "营养目标进度",
            fontSize = 16.sp,
            fontWeight = FontWeight.SemiBold,
            color = MaterialTheme.colorScheme.onSurface,
            modifier = Modifier.padding(bottom = 8.dp)
        )
        
        // 卡路里进度条
        NutrientProgressBar(
            label = "卡路里",
            progress = caloriesProgress,
            percentage = caloriesPercentage,
            current = todayCalories.toString(),
            target = targetCalories.toInt().toString(),
            color = MaterialTheme.colorScheme.primary
        )
        
        Spacer(modifier = Modifier.height(8.dp))
        
        // 蛋白质进度条
        NutrientProgressBar(
            label = "蛋白质",
            progress = proteinProgress,
            percentage = proteinPercentage,
            current = String.format(Locale.CHINA,"%.1f", todayProtein),
            target = targetProtein.toInt().toString(),
            color = MaterialTheme.colorScheme.secondary,
            unit = "g"
        )
        
        Spacer(modifier = Modifier.height(8.dp))
        
        // 脂肪进度条
        NutrientProgressBar(
            label = "脂肪",
            progress = fatProgress,
            percentage = fatPercentage,
            current = String.format(Locale.CHINA,"%.1f", todayFat),
            target = targetFat.toInt().toString(),
            color = MaterialTheme.colorScheme.tertiary,
            unit = "g"
        )
        
        Spacer(modifier = Modifier.height(8.dp))
        
        // 碳水进度条
        NutrientProgressBar(
            label = "碳水",
            progress = carbsProgress,
            percentage = carbsPercentage,
            current = String.format(Locale.CHINA,"%.1f", todayCarbs),
            target = targetCarbs.toInt().toString(),
            color = MaterialTheme.colorScheme.error,
            unit = "g"
        )

        Spacer(modifier = Modifier.height(24.dp))
        
        // 最近记录列表
        RecentRecordsList()
    }
}



/**
 * 营养素进度条组件
 * @param label 营养素标签
 * @param progress 进度值（0-1）
 * @param percentage 百分比值
 * @param current 当前值
 * @param target 目标值
 * @param color 进度条颜色
 * @param unit 单位（默认为空）
 */
@Composable
private fun NutrientProgressBar(
    label: String,
    progress: Float,
    percentage: Int,
    current: String,
    target: String,
    color: Color,
    unit: String = ""
) {
    Column(modifier = Modifier.fillMaxWidth()) {
        // 标签和百分比
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = label,
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = MaterialTheme.colorScheme.onSurface
            )
            Text(
                text = "$percentage%",
                fontSize = 14.sp,
                color = color,
                fontWeight = FontWeight.Bold
            )
        }
        
        Spacer(modifier = Modifier.height(4.dp))
        
        // 进度条
        LinearProgressIndicator(
            progress = { progress },
            modifier = Modifier
                .fillMaxWidth()
                .height(8.dp)
                .clip(CircleShape),
            color = color,
            trackColor = MaterialTheme.colorScheme.surfaceVariant,
        )
        
        Spacer(modifier = Modifier.height(4.dp))
        
        // 当前值和目标值
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Text(
                text = "$current${unit}",
                fontSize = 12.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Text(
                text = "目标: $target${unit}",
                fontSize = 12.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
        }
    }
}