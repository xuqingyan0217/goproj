package com.xqy.nutrition.data.auth

import com.coze.openapi.service.auth.JWTOAuthClient
import android.content.Context
import android.content.SharedPreferences
import android.util.Log
import java.io.IOException
import java.nio.charset.StandardCharsets
import java.util.Date
import kotlinx.coroutines.*

/**
 * Token管理对象，提供Token的获取、刷新和存储功能
 */
object Token {
    // 定义上下文和SharedPreferences相关常量
    private lateinit var applicationContext: Context
    private lateinit var sharedPreferences: SharedPreferences
    private var accessToken: String = ""
    private var tokenExpirationTime: Long = 0
    private var sessionName: String = ""
    private const val TOKEN_PREFS = "token_prefs"
    private const val KEY_ACCESS_TOKEN = "access_token"
    private const val KEY_EXPIRATION_TIME = "expiration_time"
    private const val TOKEN_REFRESH_THRESHOLD = 60 * 60 * 1000 // 1小时的毫秒数
    private const val SESSION_NAME = "session_name"

    /**
     * 初始化Token管理器
     * @param context 应用程序上下文
     */
    fun init(context: Context) {
        applicationContext = context.applicationContext
        sharedPreferences = applicationContext.getSharedPreferences(TOKEN_PREFS, Context.MODE_PRIVATE)
        loadTokenFromPrefs()
    }

    /**
     * 从SharedPreferences中加载Token和过期时间
     */
    private fun loadTokenFromPrefs() {
        accessToken = sharedPreferences.getString(KEY_ACCESS_TOKEN, "") ?: ""
        tokenExpirationTime = sharedPreferences.getLong(KEY_EXPIRATION_TIME, 0)
        sessionName = sharedPreferences.getString(SESSION_NAME, "") ?: ""
        Log.d("Token", "加载Token - accessToken: $accessToken, 过期时间: ${Date(tokenExpirationTime)}, 会话ID: $sessionName")
    }

    /**
     * 将Token和过期时间保存到SharedPreferences中
     * @param token 新的access token
     * @param expirationTime 新的过期时间
     */
    private fun saveTokenToPrefs(token: String, expirationTime: Long, sessionName: String) {
        sharedPreferences.edit().apply {
            putString(KEY_ACCESS_TOKEN, token)
            putLong(KEY_EXPIRATION_TIME, expirationTime)
            putString(SESSION_NAME, sessionName)
            apply()
        }
    }

    /**
     * 获取当前有效的AccessToken，如果Token无效则尝试刷新
     * @return 当前有效的AccessToken
     */
    suspend fun getAccessToken(): String {
        if (!isTokenValid()) {
            withContext(Dispatchers.IO) {
                refreshAccessToken()
            }
        }
        return accessToken
    }

    /**
     * 检查当前Token是否有效
     * @return 如果Token有效返回true，否则返回false
     */
    private fun isTokenValid(): Boolean {
        // 从SharedPreferences中重新加载Token
        loadTokenFromPrefs()
        return try {
            if (accessToken.isEmpty()) {
                Log.d("Token", "Token为空")
                return false
            }
            // 提前1h检查是否过期
            val currentTime = System.currentTimeMillis()
            if (currentTime >= tokenExpirationTime - TOKEN_REFRESH_THRESHOLD) {
                Log.d("Token", "Token即将过期")
                return false
            }
            true
        } catch (e: Exception) {
            Log.e("Token", "Token验证失败", e)
            false
        }
    }

    /**
     * 同步方法，用于刷新AccessToken
     */
    @Synchronized
    private fun refreshAccessToken() {
        val envConfig = EnvConfig.getInstance()
        Log.d("Token", "开始刷新Token - 获取环境配置")

        val cozeAPIBase = envConfig.cozeApiBase
        val jwtOauthClientID = envConfig.jwtOauthClientId
        val jwtOauthPrivateKeyFilePath = envConfig.jwtOauthPrivateKeyFilePath
        val jwtOauthPublicKeyID = envConfig.jwtOauthPublicKeyId

        Log.d("Token", "环境配置参数: cozeAPIBase=$cozeAPIBase, clientID=$jwtOauthClientID, keyPath=$jwtOauthPrivateKeyFilePath, publicKeyID=$jwtOauthPublicKeyID")

        val jwtOauthPrivateKey = try {
            Log.d("Token", "开始读取私钥文件: $jwtOauthPrivateKeyFilePath")
            applicationContext.assets.open(jwtOauthPrivateKeyFilePath).use { inputStream ->
                inputStream.readBytes().toString(StandardCharsets.UTF_8)
            }.also { key ->
                Log.d("Token", "私钥读取成功，长度: ${key.length}")
            }
        } catch (e: IOException) {
            Log.e("Token", "私钥读取失败", e)
            e.printStackTrace()
            return
        }

        val oauth = try {
            Log.d("Token", "开始创建OAuth客户端")
            JWTOAuthClient.JWTOAuthBuilder()
                .clientID(jwtOauthClientID)
                .privateKey(jwtOauthPrivateKey)
                .publicKey(jwtOauthPublicKeyID)
                .baseURL(cozeAPIBase)
                .build()
                .also { Log.d("Token", "OAuth客户端创建成功") }
        } catch (e: Exception) {
            Log.e("Token", "OAuth客户端创建失败", e)
            e.printStackTrace()
            return
        }

        try {
            Log.d("Token", "开始获取access token")
            // 添加会话，便于网页组件使用
            val tokenResponse = oauth.getAccessToken(86400,"10001")
            accessToken = tokenResponse.accessToken
            Log.d("Token", "获取access token成功: ${accessToken.take(10)}...")
            // 设置token过期时间为24小时 86400s
            tokenExpirationTime = System.currentTimeMillis() + (24 * 60 * 60 * 1000)
            sessionName = "10010"
            saveTokenToPrefs(accessToken, tokenExpirationTime, sessionName)
            Log.d("Token", "Token刷新完成，已保存到SharedPreferences")
        } catch (e: Exception) {
            Log.e("Token", "获取access token失败", e)
            e.printStackTrace()
            throw e
        }
    }
}
