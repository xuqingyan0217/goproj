package com.xqy.nutrition.data.model

import android.content.Context
import android.content.SharedPreferences
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.net.Uri
import android.util.Base64
import android.util.Log
import java.io.ByteArrayOutputStream

/**
 * 用户个人资料管理类
 * 负责存储和检索用户的个人信息，如头像、用户名和目标
 */
class UserProfileManager private constructor(context: Context) {
    
    // SharedPreferences实例，用于数据持久化
    private val sharedPreferences: SharedPreferences =
        context.getSharedPreferences(PREF_NAME, Context.MODE_PRIVATE)

    companion object {
        private const val TAG = "UserProfileManager"
        private const val PREF_NAME = "user_profile_prefs"
        private const val KEY_USERNAME = "username"
        private const val KEY_AVATAR = "avatar"
        private const val KEY_WEIGHT_GOAL = "weight_goal"
        private const val KEY_CALORIES = "calories"
        private const val KEY_PROTEIN = "protein"
        private const val KEY_CARBS = "carbs"
        private const val KEY_FAT = "fat"
        private const val DEFAULT_USERNAME = "UserT"
        private const val DEFAULT_WEIGHT_GOAL = "-5"
        private const val DEFAULT_CALORIES = "2500"
        private const val DEFAULT_PROTEIN = "90"
        private const val DEFAULT_CARBS = "225"
        private const val DEFAULT_FAT = "60"
        
        @Volatile
        private var instance: UserProfileManager? = null
        
        fun getInstance(context: Context): UserProfileManager {
            return instance ?: synchronized(this) {
                instance ?: UserProfileManager(context.applicationContext).also { instance = it }
            }
        }
    }

    /**
     * 获取用户名
     * @return 用户名，如果未设置则返回默认值
     */
    fun getUsername(): String {
        return sharedPreferences.getString(KEY_USERNAME, DEFAULT_USERNAME) ?: DEFAULT_USERNAME
    }
    
    /**
     * 设置用户名
     * @param username 新的用户名
     */
    fun setUsername(username: String) {
        sharedPreferences.edit().putString(KEY_USERNAME, username).apply()
        Log.d(TAG, "Username updated to: $username")
    }
    
    /**
     * 获取减重目标
     * @return 减重目标值，如果未设置则返回默认值
     */
    fun getWeightGoal(): String {
        return sharedPreferences.getString(KEY_WEIGHT_GOAL, DEFAULT_WEIGHT_GOAL) ?: DEFAULT_WEIGHT_GOAL
    }
    
    /**
     * 设置减重目标
     * @param goal 新的减重目标值
     */
    fun setWeightGoal(goal: String) {
        sharedPreferences.edit().putString(KEY_WEIGHT_GOAL, goal).apply()
        Log.d(TAG, "Weight goal updated to: $goal")
    }
    
    /**
     * 获取每日热量目标
     * @return 每日热量目标值，如果未设置则返回默认值
     */
    fun getCalories(): String {
        return sharedPreferences.getString(KEY_CALORIES, DEFAULT_CALORIES) ?: DEFAULT_CALORIES
    }
    
    /**
     * 设置每日热量目标
     * @param calories 新的每日热量目标值
     */
    fun setCalories(calories: String) {
        sharedPreferences.edit().putString(KEY_CALORIES, calories).apply()
        Log.d(TAG, "Calories goal updated to: $calories")
    }
    
    /**
     * 获取蛋白质目标
     * @return 蛋白质目标值，如果未设置则返回默认值
     */
    fun getProtein(): String {
        return sharedPreferences.getString(KEY_PROTEIN, DEFAULT_PROTEIN) ?: DEFAULT_PROTEIN
    }
    
    /**
     * 设置蛋白质目标
     * @param protein 新的蛋白质目标值
     */
    fun setProtein(protein: String) {
        sharedPreferences.edit().putString(KEY_PROTEIN, protein).apply()
        Log.d(TAG, "Protein goal updated to: $protein")
    }
    
    /**
     * 获取碳水化合物目标
     * @return 碳水化合物目标值，如果未设置则返回默认值
     */
    fun getCarbs(): String {
        return sharedPreferences.getString(KEY_CARBS, DEFAULT_CARBS) ?: DEFAULT_CARBS
    }
    
    /**
     * 设置碳水化合物目标
     * @param carbs 新的碳水化合物目标值
     */
    fun setCarbs(carbs: String) {
        sharedPreferences.edit().putString(KEY_CARBS, carbs).apply()
        Log.d(TAG, "Carbs goal updated to: $carbs")
    }
    
    /**
     * 获取脂肪目标
     * @return 脂肪目标值，如果未设置则返回默认值
     */
    fun getFat(): String {
        return sharedPreferences.getString(KEY_FAT, DEFAULT_FAT) ?: DEFAULT_FAT
    }
    
    /**
     * 设置脂肪目标
     * @param fat 新的脂肪目标值
     */
    fun setFat(fat: String) {
        sharedPreferences.edit().putString(KEY_FAT, fat).apply()
        Log.d(TAG, "Fat goal updated to: $fat")
    }
    
    /**
     * 保存头像
     * @param bitmap 头像位图
     */
    private fun saveAvatar(bitmap: Bitmap) {
        val base64Avatar = bitmapToBase64(bitmap)
        sharedPreferences.edit().putString(KEY_AVATAR, base64Avatar).apply()
        Log.d(TAG, "Avatar updated")
    }
    
    /**
     * 保存头像
     * @param uri 头像URI
     * @param context 上下文
     */
    fun saveAvatar(uri: Uri, context: Context) {
        try {
            val inputStream = context.contentResolver.openInputStream(uri)
            val bitmap = BitmapFactory.decodeStream(inputStream)
            inputStream?.close()
            
            if (bitmap != null) {
                saveAvatar(bitmap)
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error saving avatar from URI", e)
        }
    }
    
    /**
     * 获取头像
     * @return 头像位图，如果未设置则返回null
     */
    fun getAvatar(): Bitmap? {
        val base64Avatar = sharedPreferences.getString(KEY_AVATAR, null) ?: return null
        return try {
            base64ToBitmap(base64Avatar)
        } catch (e: Exception) {
            Log.e(TAG, "Error loading avatar", e)
            null
        }
    }
    

    
    /**
     * 将Bitmap转换为Base64编码的字符串
     * @param bitmap 位图
     * @return Base64编码的字符串
     */
    private fun bitmapToBase64(bitmap: Bitmap): String {
        val byteArrayOutputStream = ByteArrayOutputStream()
        // 压缩图片以减小大小，同时保持足够的质量
        bitmap.compress(Bitmap.CompressFormat.JPEG, 70, byteArrayOutputStream)
        val byteArray = byteArrayOutputStream.toByteArray()
        return Base64.encodeToString(byteArray, Base64.DEFAULT)
    }
    
    /**
     * 将Base64编码的字符串转换为Bitmap
     * @param base64String Base64编码的字符串
     * @return 位图
     */
    private fun base64ToBitmap(base64String: String): Bitmap {
        val decodedBytes = Base64.decode(base64String, Base64.DEFAULT)
        return BitmapFactory.decodeByteArray(decodedBytes, 0, decodedBytes.size)
    }
}