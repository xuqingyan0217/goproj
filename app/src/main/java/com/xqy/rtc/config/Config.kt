package com.xqy.rtc.config

import android.content.Context
import com.xqy.nutrition.data.auth.EnvConfig
import com.xqy.nutrition.data.auth.Token
import java.lang.ref.WeakReference

/**
 * 单例配置类，用于管理应用中需要的配置信息
 */
class Config private constructor() {
    // 环境配置实例
    private val envConfig = EnvConfig.getInstance()

    /**
     * 悬挂函数，用于获取Coze访问令牌
     *
     * @return String 类型的访问令牌
     */
    suspend fun getCozeAccessToken(): String = Token.getAccessToken()

    // API基础URL，通过环境配置获取，并且只能在类内部设置
    var baseURL: String = envConfig.cozeApiBase
        private set

    // 机器人ID，通过环境配置获取，并且只能在类内部设置
    var botID: String = envConfig.botId
        private set

    // 语音ID，通过环境配置获取，并且只能在类内部设置
    var voiceID: String = envConfig.voiceId
        private set

    // 伴生对象，用于实现单例模式
    companion object {
        // 单例实例
        @Volatile
        private var instance: Config? = null
        
        // 使用WeakReference存储上下文，避免内存泄漏
        private var contextRef: WeakReference<Context>? = null

        /**
         * 初始化配置
         *
         * @param appContext 应用上下文
         */
        fun init(appContext: Context) {
            contextRef = WeakReference(appContext.applicationContext)
            // 先初始化EnvConfig
            EnvConfig.init(appContext)
            // 再初始化Token
            Token.init(appContext)
            getInstance()
        }

        /**
         * 同步方法，用于初始化单例实例
         *
         * @return Config 类的单例实例
         */
        @Synchronized
        private fun initInstance(): Config {
            if (instance == null) {
                instance = Config()
            }
            return instance!!
        }

        /**
         * 获取Config的单例实例
         *
         * @return Config 类的单例实例
         */
        fun getInstance(): Config {
            val context = contextRef?.get()
            checkNotNull(context) { "Config must be initialized with context first" }
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
