package com.xqy.nutrition.data.network

import android.annotation.SuppressLint
import com.xqy.nutrition.data.auth.EnvConfig
import com.xqy.nutrition.data.model.FoodNutrition
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONArray
import org.json.JSONObject
import java.util.concurrent.TimeUnit

/**
 * 火山引擎API服务
 */
object ApiService {
    @SuppressLint("StaticFieldLeak")
    private val envConfig = EnvConfig.getInstance()
    private val API_KEY = envConfig.apiKey
    private val BASE_URL = envConfig.baseUrl
    private val MODEL = envConfig.dbModel
    
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()
    
    /**
     * 发送图片识别请求
     * @param base64Image Base64编码的图片
     * @return 识别结果
     */
    fun recognizeFood(base64Image: String): Result<FoodNutrition> {
        return try {
            // 构建请求体
            val requestBody = createRequestBody(base64Image)
            
            // 构建请求
            val request = Request.Builder()
                .url(BASE_URL)
                .addHeader("Content-Type", "application/json")
                .addHeader("Authorization", "Bearer $API_KEY")
                .post(requestBody)
                .build()
            
            // 执行请求
            val response = client.newCall(request).execute()
            val responseBody = response.body?.string()
            
            if (response.isSuccessful && responseBody != null) {
                // 解析响应
                val nutrition = parseResponse(responseBody)
                Result.success(nutrition)
            } else {
                Result.failure(Exception("API请求失败: ${response.code}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    /**
     * 创建请求体
     */
    private fun createRequestBody(base64Image: String): okhttp3.RequestBody {
        val jsonObject = JSONObject()
        jsonObject.put("model", MODEL)
        
        val messagesArray = JSONArray()
        val messageObject = JSONObject()
        messageObject.put("role", "user")
        
        val contentArray = JSONArray()
        
        // 添加图片
        val imageObject = JSONObject()
        imageObject.put("type", "image_url")
        
        val imageUrlObject = JSONObject()
        imageUrlObject.put("url", "data:image/jpeg;base64,$base64Image")
        
        imageObject.put("image_url", imageUrlObject)
        contentArray.put(imageObject)
        
        // 添加文本提示
        val textObject = JSONObject()
        textObject.put("type", "text")
        textObject.put("text", "请识别这张图片中的食物，并提供以下营养信息：1. 食物名称 2. 卡路里(kcal) 3. 蛋白质(g) 4. 脂肪(g) 5. 碳水化合物(g)。请只返回JSON格式数据，格式为：{\"name\":\"食物名称\",\"calories\":\"卡路里值\",\"protein\":\"蛋白质值\",\"fat\":\"脂肪值\",\"carbs\":\"碳水化合物值\"}")
        contentArray.put(textObject)
        
        messageObject.put("content", contentArray)
        messagesArray.put(messageObject)
        
        jsonObject.put("messages", messagesArray)
        
        return jsonObject.toString().toRequestBody("application/json".toMediaTypeOrNull())
    }
    
    /**
     * 解析API响应
     */
    private fun parseResponse(responseBody: String): FoodNutrition {
        try {
            android.util.Log.d("ApiService", "开始解析响应: $responseBody")
            val jsonResponse = JSONObject(responseBody)
            val choices = jsonResponse.getJSONArray("choices")
            if (choices.length() > 0) {
                val firstChoice = choices.getJSONObject(0)
                val message = firstChoice.getJSONObject("message")
                val content = message.getString("content")
                
                android.util.Log.d("ApiService", "提取到的内容: $content")
                
                // 从内容中提取JSON
                val jsonStart = content.indexOf("{")
                val jsonEnd = content.lastIndexOf("}") + 1
                
                if (jsonStart in 0 until jsonEnd) {
                    val nutritionJson = content.substring(jsonStart, jsonEnd)
                    android.util.Log.d("ApiService", "提取到的JSON: $nutritionJson")
                    
                    val nutritionObject = JSONObject(nutritionJson)
                    
                    // 提取各个营养素值并记录日志
                    val name = nutritionObject.optString("name", "未知食物")
                    val caloriesStr = nutritionObject.optString("calories", "0")
                    val proteinStr = nutritionObject.optString("protein", "0g")
                    val fatStr = nutritionObject.optString("fat", "0g")
                    val carbsStr = nutritionObject.optString("carbs", "0g")
                    
                    android.util.Log.d("ApiService", "原始值 - 名称: $name, 卡路里: $caloriesStr, 蛋白质: $proteinStr, 脂肪: $fatStr, 碳水: $carbsStr")
                    
                    // 处理数值
                    val calories = caloriesStr.replace("[^0-9.]".toRegex(), "").toFloatOrNull() ?: 0f
                    
                    // 改进蛋白质、脂肪和碳水值的提取方法
                    val protein = proteinStr.replace("[^0-9.]".toRegex(), "").toFloatOrNull() ?: 0f
                    val fat = fatStr.replace("[^0-9.]".toRegex(), "").toFloatOrNull() ?: 0f
                    val carbs = carbsStr.replace("[^0-9.]".toRegex(), "").toFloatOrNull() ?: 0f
                    
                    android.util.Log.d("ApiService", "处理后的值 - 卡路里: $calories, 蛋白质: $protein, 脂肪: $fat, 碳水: $carbs")
                    
                    // 确保蛋白质、脂肪和碳水值不为0，如果为0则使用默认值
                    val finalProtein = if (protein <= 0f) 0.1f else protein
                    val finalFat = if (fat <= 0f) 0.1f else fat
                    val finalCarbs = if (carbs <= 0f) calories * 0.15f / 4f else carbs // 如果API没返回碳水，则估算
                    
                    android.util.Log.d("ApiService", "最终返回值 - 卡路里: $calories, 蛋白质: $finalProtein, 脂肪: $finalFat, 碳水: $finalCarbs")
                    
                    return FoodNutrition(
                        name = name,
                        calories = calories,
                        protein = finalProtein,
                        fat = finalFat,
                        carbs = finalCarbs
                    )
                } else {
                    android.util.Log.w("ApiService", "无法从内容中提取JSON, jsonStart: $jsonStart, jsonEnd: $jsonEnd")
                }
            } else {
                android.util.Log.w("ApiService", "API响应中没有choices")
            }
            
            // 默认返回
            android.util.Log.w("ApiService", "使用默认值返回")
            return FoodNutrition("未能识别", 0f, 0f, 0f, 0f)
        } catch (e: Exception) {
            android.util.Log.e("ApiService", "解析响应时发生异常", e)
            e.printStackTrace()
            return FoodNutrition("解析错误", 0f, 0f, 0f, 0f)
        }
    }
}