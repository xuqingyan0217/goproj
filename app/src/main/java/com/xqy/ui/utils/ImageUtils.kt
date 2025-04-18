package com.xqy.ui.utils

import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.net.Uri
import android.util.Base64
import java.io.ByteArrayOutputStream
import java.io.IOException

/**
 * 图片处理工具类
 */
object ImageUtils {
    
    /**
     * 将Uri转换为Base64编码的字符串
     * @param context 上下文
     * @param uri 图片Uri
     * @return Base64编码的字符串，失败返回null
     */
    fun uriToBase64(context: Context, uri: Uri): String? {
        return try {
            val inputStream = context.contentResolver.openInputStream(uri)
            val bitmap = BitmapFactory.decodeStream(inputStream)
            inputStream?.close()
            
            if (bitmap != null) {
                bitmapToBase64(bitmap)
            } else {
                null
            }
        } catch (e: IOException) {
            e.printStackTrace()
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
        // 压缩图片以减小大小，同时保持足够的质量用于识别
        bitmap.compress(Bitmap.CompressFormat.JPEG, 80, byteArrayOutputStream)
        val byteArray = byteArrayOutputStream.toByteArray()
        return Base64.encodeToString(byteArray, Base64.NO_WRAP)
    }
}