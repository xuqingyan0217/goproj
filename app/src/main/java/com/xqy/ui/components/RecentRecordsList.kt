/**
 * 最近饮食记录列表组件模块
 * 用于显示用户的饮食记录列表，包括食物名称、时间、卡路里和营养成分信息
 */
package com.xqy.ui.components

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import java.util.Locale

/**
 * 食物记录数据模型
 * @property name 食物名称
 * @property time 记录时间
 * @property calories 卡路里
 * @property protein 蛋白质含量（克）
 * @property fat 脂肪含量（克）
 * @property carbs 碳水化合物含量（克）
 * @property date 记录日期
 */
data class FoodRecord(
    val name: String,
    val time: String,
    val calories: Float,
    val protein: Float = 0f,
    val fat: Float = 0f,
    val carbs: Float = 0f,
    val date: String = ""
)

/**
 * 最近饮食记录列表组件
 * 显示当前选定日期的所有饮食记录
 */
@Composable
fun RecentRecordsList() {
    // 从FoodRecordManager获取当前选定日期的记录数据
    val records = remember(com.xqy.nutrition.data.model.FoodRecordManager.records) {
        com.xqy.nutrition.data.model.FoodRecordManager.records
    }
    
    // 获取当前选定的日期
    val selectedDate = remember {
        mutableStateOf(com.xqy.nutrition.data.model.FoodRecordManager.formatDateForDisplay(com.xqy.nutrition.data.model.FoodRecordManager.getSelectedDate()))
    }
    
    // 当记录数据变化时更新选定日期
    LaunchedEffect(com.xqy.nutrition.data.model.FoodRecordManager.records) {
        selectedDate.value = com.xqy.nutrition.data.model.FoodRecordManager.formatDateForDisplay(com.xqy.nutrition.data.model.FoodRecordManager.getSelectedDate())
    }

    Column(modifier = Modifier.fillMaxWidth()) {
        // 日期标题
        Text(
            text = "${selectedDate.value}记录",
            fontSize = 20.sp,
            fontWeight = FontWeight.Medium,
            modifier = Modifier.padding(bottom = 16.dp)
        )

        // 记录列表 - 使用Column代替LazyColumn以避免嵌套滚动冲突
        Column(
            modifier = Modifier.fillMaxWidth(),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            records.forEach { record ->
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.surface
                    )
                ) {
                    RecordItem(record)
                }
            }
        }
    }
}

/**
 * 单条记录项组件
 * @param record 食物记录数据
 */
@Composable
fun RecordItem(record: FoodRecord) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surface
        )
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            // 食物名称和时间行
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                // 左侧：食物图标和名称
                Row(verticalAlignment = Alignment.CenterVertically) {
                    // 食物首字母图标
                    Surface(
                        modifier = Modifier.size(48.dp),
                        shape = MaterialTheme.shapes.medium,
                        color = MaterialTheme.colorScheme.primaryContainer
                    ) {
                        Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                            Text(
                                text = record.name.take(1).uppercase(),
                                fontSize = 20.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.onPrimaryContainer
                            )
                        }
                    }
                    
                    Spacer(modifier = Modifier.width(16.dp))
                    
                    // 食物名称和时间
                    Column {
                        Text(
                            text = record.name,
                            fontSize = 18.sp,
                            fontWeight = FontWeight.SemiBold
                        )
                        Text(
                            text = record.time,
                            fontSize = 14.sp,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                }
                
                // 右侧：删除按钮
                IconButton(
                    onClick = {
                        com.xqy.nutrition.data.model.FoodRecordManager.removeRecord(record)
                    }
                ) {
                    Icon(
                        imageVector = androidx.compose.material.icons.Icons.Filled.Delete,
                        contentDescription = "删除",
                        tint = MaterialTheme.colorScheme.error
                    )
                }
            }
            
            // 营养素信息行
            Spacer(modifier = Modifier.height(12.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                NutrientInfo("卡路里", record.calories, MaterialTheme.colorScheme.primary)
                NutrientInfo("蛋白质", record.protein, MaterialTheme.colorScheme.secondary)
                NutrientInfo("脂肪", record.fat, MaterialTheme.colorScheme.tertiary)
                NutrientInfo("碳水", record.carbs, MaterialTheme.colorScheme.error)
            }
        }
    }
}

/**
 * 营养素信息显示组件
 * @param label 营养素标签
 * @param value 营养素含量
 * @param color 显示颜色
 */
@Composable
private fun NutrientInfo(label: String, value: Float, color: Color) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(
            text = "${String.format(Locale.CHINA, "%.1f", value)}g",
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