package com.xqy.nutrition.data.auth

import android.content.Context
import com.xqy.nutrition.R
import java.lang.ref.WeakReference

/**
 * 环境配置类，用于获取应用中使用的各种API和认证相关的配置信息
 */
class EnvConfig private constructor() {

    /**
     * 获取Coze API的基础URL
     */
    val cozeApiBase: String
        get() = contextRef?.get()?.getString(R.string.COZE_API_BASE) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取JWT认证所需的客户端ID
     */
    val jwtOauthClientId: String
        get() = contextRef?.get()?.getString(R.string.COZE_JWT_OAUTH_CLIENT_ID) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取JWT认证所需的私钥文件路径
     */
    val jwtOauthPrivateKeyFilePath: String
        get() = contextRef?.get()?.getString(R.string.COZE_JWT_OAUTH_PRIVATE_KEY_FILE_PATH) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取JWT认证所需的公钥ID
     */
    val jwtOauthPublicKeyId: String
        get() = contextRef?.get()?.getString(R.string.COZE_JWT_OAUTH_PUBLIC_KEY_ID) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取机器人ID
     */
    val botId : String
        get() = contextRef?.get()?.getString(R.string.BOT_ID) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取语音ID
     */
    val voiceId : String
        get() = contextRef?.get()?.getString(R.string.VOICE_ID) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取API密钥
     */
    val apiKey : String
        get() = contextRef?.get()?.getString(R.string.API_KEY) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取基础URL
     */
    val baseUrl : String
        get() = contextRef?.get()?.getString(R.string.BASE_URL) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 获取数据库模型
     */
    val dbModel : String
        get() = contextRef?.get()?.getString(R.string.MODEL) ?: throw IllegalStateException("Context is not initialized")

    /**
     * 伴生对象，用于实现单例模式并初始化配置
     */
    companion object {
        private var instance: EnvConfig? = null
        private var contextRef: WeakReference<Context>? = null

        /**
         * 初始化配置，必须在应用初始化时调用
         * @param context 应用上下文，用于获取资源字符串
         */
        fun init(context: Context) {
            contextRef = WeakReference(context.applicationContext)
        }

        /**
         * 同步初始化实例，确保在多线程环境下安全地创建单例
         * @return 初始化的EnvConfig实例
         */
        @Synchronized
        private fun initInstance(): EnvConfig {
            if (instance == null) {
                instance = EnvConfig()
            }
            return instance!!
        }

        /**
         * 获取EnvConfig的实例，如果未初始化则会先进行初始化
         * @return EnvConfig的实例
         */
        fun getInstance(): EnvConfig {
            val context = contextRef?.get()
            checkNotNull(context) { "EnvConfig must be initialized with context first" }
            return initInstance()
        }

        /**
         * 清理资源，防止内存泄漏
         */
        fun cleanup() {
            instance = null
            contextRef = null
        }
    }
}
