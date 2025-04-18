package com.xqy.nutrition.data.model

import android.content.Context
import android.content.SharedPreferences
import android.util.Log
import androidx.compose.runtime.mutableStateListOf
import com.xqy.ui.components.FoodRecord
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.util.Locale

/**
 * 食物记录管理器
 * 负责管理用户的食物摄入记录，包括添加、删除、保存、加载和按日期筛选记录等功能。
 * 使用SharedPreferences持久化存储数据。
 */
object FoodRecordManager {
    // SharedPreferences实例，用于数据持久化
    private lateinit var sharedPreferences: SharedPreferences
    // SharedPreferences文件名和键值
    private const val PREF_NAME = "food_records"
    private const val RECORDS_KEY = "records_list"
    
    // 当前选择日期的记录列表
    val records = mutableStateListOf<FoodRecord>()
    // 所有日期的记录列表
    val allRecords = mutableStateListOf<FoodRecord>()
    
    // 当前选择的日期，默认为今天（使用中国时区）
    private var selectedDate = LocalDate.now(ZoneId.of("Asia/Shanghai"))
    
    /**
     * 初始化FoodRecordManager
     * 加载已保存的记录并根据当前选择的日期筛选记录
     * @param context 应用程序上下文，用于获取SharedPreferences
     */
    fun initialize(context: Context) {
        Log.d("FoodRecordManager", "Initializing with context: $context")
        sharedPreferences = context.getSharedPreferences(PREF_NAME, Context.MODE_PRIVATE)
        loadRecords()
        
        // 添加昨天的测试数据
        /*val yesterday = LocalDate.now(ZoneId.of("Asia/Shanghai")).minusDays(1)
        val yesterdayStr = yesterday.format(DateTimeFormatter.ofPattern("yyyy-MM-dd"))
        
        // 添加昨天的早餐记录
        allRecords.add(FoodRecord("早餐-牛奶", "08:30", 150f, 8.0f, 3.5f, 5.0f, yesterdayStr))
        allRecords.add(FoodRecord("早餐-面包", "08:30", 200f, 6.0f, 4.0f, 5.0f, yesterdayStr))
        
        // 添加昨天的午餐记录
        allRecords.add(FoodRecord("午餐-米饭", "12:30", 250f, 5.0f, 0.5f, 5.0f, yesterdayStr))
        allRecords.add(FoodRecord("午餐-鱼香肉丝", "12:30", 300f, 15.0f, 12.0f, 5.0f, yesterdayStr))
        
        // 添加昨天的晚餐记录
        allRecords.add(FoodRecord("晚餐-水饺", "18:30", 400f, 12.0f, 8.0f, 5.0f, yesterdayStr))
        allRecords.add(FoodRecord("晚餐-蔬菜沙拉", "18:30", 100f, 2.0f, 3.0f, 5.0f, yesterdayStr))
        
        saveRecords()*/
        filterRecordsByDate(selectedDate)
        Log.d("FoodRecordManager", "Initialization complete, loaded ${records.size} records with test data")
    }
    
    /**
     * 从SharedPreferences加载保存的食物记录
     * 解析记录字符串并转换为FoodRecord对象
     */
    private fun loadRecords() {
        val recordsStr = sharedPreferences.getString(RECORDS_KEY, "") ?: ""
        Log.d("FoodRecordManager", "Loading records, raw data: $recordsStr")
        
        if (recordsStr.isNotEmpty()) {
            allRecords.clear()
            val recordsList = recordsStr.split("|") // 使用|分隔不同记录
            Log.d("FoodRecordManager", "Found ${recordsList.size} total records")
            
            recordsList.forEach { recordStr ->
                val parts = recordStr.split(",") // 使用,分隔记录的各个字段
                if (parts.size >= 7) { // 确保有足够的部分包含日期
                    try {
                        // 解析蛋白质、脂肪和碳水值
                        val proteinStr = parts[3]
                        val fatStr = parts[4]
                        val carbsStr = parts[5] // 碳水在第6个位置 (index 5)
                        val proteinValue = proteinStr.replace("g", "").toFloatOrNull() ?: 0f
                        val fatValue = fatStr.replace("g", "").toFloatOrNull() ?: 0f
                        val carbsValue = carbsStr.replace("g", "").toFloatOrNull() ?: 0f
                        
                        Log.d("FoodRecordManager", "解析记录 - 原始蛋白质: $proteinStr, 解析后: $proteinValue, 原始脂肪: $fatStr, 解析后: $fatValue, 原始碳水: $carbsStr, 解析后: $carbsValue")
                        
                        // 处理时间格式
                        val timeStr = parts[1]
                        Log.d("FoodRecordManager", "原始时间字符串: $timeStr")
                        
                        val formattedTime = try {
                            if (timeStr.matches(Regex("\\d{2}:\\d{2}"))) {
                                timeStr
                            } else {
                                LocalDateTime.now(ZoneId.of("Asia/Shanghai")).format(DateTimeFormatter.ofPattern("HH:mm"))
                            }
                        } catch (e: Exception) {
                            Log.e("FoodRecordManager", "时间格式解析错误: $timeStr", e)
                            LocalDateTime.now(ZoneId.of("Asia/Shanghai")).format(DateTimeFormatter.ofPattern("HH:mm"))
                        }
                        
                        Log.d("FoodRecordManager", "处理后的时间: $formattedTime")
                        
                        // 获取日期 - 日期在第7个位置 (index 6)
                        val dateStr = parts[6]
                        Log.d("FoodRecordManager", "日期字符串: $dateStr")
                        
                        // 创建并添加记录
                        val record = FoodRecord(
                            name = parts[0],
                            time = formattedTime,
                            calories = parts[2].toFloat(),
                            protein = proteinValue,
                            fat = fatValue,
                            carbs = carbsValue,
                            date = dateStr
                        )
                        
                        // 尝试解析日期以验证格式是否正确
                        try {
                            LocalDate.parse(dateStr)
                            Log.d("FoodRecordManager", "日期解析成功: $dateStr")
                        } catch (e: Exception) {
                            Log.e("FoodRecordManager", "日期解析失败: $dateStr", e)
                        }
                        
                        Log.d("FoodRecordManager", "添加记录到列表 - 名称: ${record.name}, 时间: ${record.time}, 蛋白质: ${record.protein}, 脂肪: ${record.fat}, 碳水: ${record.carbs}")
                        allRecords.add(record)
                    } catch (e: Exception) {
                        Log.e("FoodRecordManager", "Error parsing record: $recordStr", e)
                    }
                } else {
                    Log.w("FoodRecordManager", "Invalid record format: $recordStr")
                }
            }
            
            // 根据当前选择的日期筛选记录
            filterRecordsByDate(selectedDate)
            
            Log.d("FoodRecordManager", "Loaded ${allRecords.size} total records, ${records.size} records for selected date")
        } else {
            Log.d("FoodRecordManager", "No records found in SharedPreferences")
        }
    }
    
    /**
     * 添加新的食物记录
     * @param name 食物名称
     * @param calories 卡路里
     * @param protein 蛋白质含量
     * @param fat 脂肪含量
     * @param carbs 碳水化合物含量
     */
    fun addRecord(name: String, calories: Float, protein: Float, fat: Float, carbs: Float = 0f) {
        Log.d("FoodRecordManager", "添加记录 - 名称: $name, 卡路里: $calories, 蛋白质: $protein, 脂肪: $fat")
        // 使用中国时区获取当前时间
        val currentTime = LocalDateTime.now(ZoneId.of("Asia/Shanghai"))
        val timeFormatter = DateTimeFormatter.ofPattern("HH:mm")
        val dateFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd")
        val currentDate = LocalDate.now(ZoneId.of("Asia/Shanghai")).format(dateFormatter)
        
        // 格式化时间
        val formattedTime = currentTime.format(timeFormatter)
        Log.d("FoodRecordManager", "当前时间: $currentTime, 格式化后: $formattedTime, 时区: ${ZoneId.of("Asia/Shanghai")}")
        
        // 创建新记录
        val record = FoodRecord(
            name = name,
            time = formattedTime,
            calories = calories,
            protein = protein,
            fat = fat,
            carbs = carbs,
            date = currentDate
        )
        
        Log.d("FoodRecordManager", "创建记录对象 - 蛋白质: ${record.protein}, 脂肪: ${record.fat}")
        allRecords.add(record)
        
        // 如果记录日期与当前选择的日期相同，添加到筛选后的记录列表
        val recordDate = LocalDate.parse(currentDate)
        if (recordDate.equals(selectedDate)) {
            records.add(record)
        }
        
        saveRecords()
    }
    
    /**
     * 将所有记录保存到SharedPreferences
     * 记录格式：名称,时间,卡路里,蛋白质,脂肪,日期
     * 不同记录之间使用|分隔
     */
    private fun saveRecords() {
        // 检查记录的营养值
        allRecords.forEachIndexed { index, record ->
            Log.d("FoodRecordManager", "保存前记录[$index] - 名称: ${record.name}, 蛋白质: ${record.protein}, 脂肪: ${record.fat}")
        }
        
        // 格式化记录为字符串
        val recordsStr = allRecords.joinToString("|") { record ->
            // 确保营养值不为0
            val proteinValue = if (record.protein <= 0f) 0.1f else record.protein
            val fatValue = if (record.fat <= 0f) 0.1f else record.fat
            val carbsValue = if (record.carbs <= 0f) 0.1f else record.carbs
            
            // 格式化营养值
            val formattedProtein = String.format(Locale.CHINA, "%.1f", proteinValue)
            val formattedFat = String.format(Locale.CHINA, "%.1f", fatValue)
            val formattedCarbs = String.format(Locale.CHINA, "%.1f", carbsValue)
            Log.d("FoodRecordManager", "格式化 - 原始蛋白质: ${record.protein}, 格式化后: ${formattedProtein}g, 原始脂肪: ${record.fat}, 格式化后: ${formattedFat}g, 原始碳水: ${record.carbs}, 格式化后: ${formattedCarbs}g")
            
            "${record.name},${record.time},${record.calories.toInt()},${formattedProtein}g,${formattedFat}g,${formattedCarbs}g,${record.date}"
        }
        
        Log.d("FoodRecordManager", "保存 ${records.size} 条记录, 数据: $recordsStr")
        sharedPreferences.edit().putString(RECORDS_KEY, recordsStr).apply()
    }
    
    /**
     * 获取选定日期的营养总量
     * @return Quadruple<总卡路里, 蛋白质, 脂肪, 碳水>
     */
    fun getTodayTotalNutrition(): Quadruple<Int, String, String, String> {
        var totalCalories = 0
        var totalProtein = 0f
        var totalFat = 0f
        var totalCarbs = 0f
        
        // 计算当前筛选记录的营养总量
        if (records.isEmpty()) {
            Log.d("FoodRecordManager", "当前日期没有记录，返回0值")
            return Quadruple(0, "0.0g", "0.0g", "0.0g")
        }
        
        records.forEach { record ->
            totalCalories += record.calories.toInt()
            totalProtein += record.protein
            totalFat += record.fat
            totalCarbs += record.carbs
        }
        
        return Quadruple(
            totalCalories,
            "${String.format(Locale.CHINA, "%.1f", totalProtein)}g",
            "${String.format(Locale.CHINA, "%.1f", totalFat)}g",
            "${String.format(Locale.CHINA, "%.1f", totalCarbs)}g"
        )
    }
    
    /**
     * 根据日期筛选记录
     * @param date 要筛选的日期
     */
    fun filterRecordsByDate(date: LocalDate) {
        Log.d("FoodRecordManager", "筛选日期: $date, 总记录数: ${allRecords.size}")
        selectedDate = date
        
        // 清空当前记录列表并重新筛选
        records.clear()
        allRecords.forEach { record ->
            try {
                val recordDate = LocalDate.parse(record.date)
                if (recordDate.equals(date)) {
                    records.add(record)
                }
            } catch (e: Exception) {
                Log.e("FoodRecordManager", "日期解析错误: ${record.date}", e)
            }
        }
        
        Log.d("FoodRecordManager", "筛选后记录数: ${records.size}")
    }
    
    /**
     * 获取当前选择的日期
     * @return 当前选择的日期
     */
    fun getSelectedDate(): LocalDate {
        return selectedDate
    }
    
    /**
     * 格式化日期为显示格式（年月日）
     * @param date 要格式化的日期
     * @return 格式化后的日期字符串
     */
    fun formatDateForDisplay(date: LocalDate): String {
        return date.format(DateTimeFormatter.ofPattern("yyyy年M月d日"))
    }
    

    /**
     * 删除指定的食物记录
     * @param record 要删除的食物记录
     */
    fun removeRecord(record: FoodRecord) {
        Log.d("FoodRecordManager", "删除记录 - 名称: ${record.name}, 时间: ${record.time}")
        records.remove(record)
        allRecords.remove(record)
        saveRecords()
        Log.d("FoodRecordManager", "删除完成，剩余记录数: ${records.size}, 总记录数: ${allRecords.size}")
    }
}