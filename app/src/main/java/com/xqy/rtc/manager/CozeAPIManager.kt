package com.xqy.rtc.manager

import com.coze.openapi.service.auth.TokenAuth
import com.coze.openapi.service.service.CozeAPI
import com.xqy.rtc.config.Config
import kotlinx.coroutines.runBlocking

/**
 * CozeAPI管理类，用于初始化和提供CozeAPI实例
 * 该类采用单例模式，确保在整个应用中只有一个实例
 */
class CozeAPIManager private constructor() {
    private var cozeAPI: CozeAPI

    init {
        // 在构造器中初始化CozeAPI实例，使用协程框架的runBlocking以确保在创建实例时完成初始化
        runBlocking {
            cozeAPI = initCozeAPI()
        }
    }

    /**
     * 初始化CozeAPI实例的 suspend 函数
     * 使用CozeAPI的Builder模式进行配置和构建
     *
     * @return CozeAPI实例
     */
    private suspend fun initCozeAPI(): CozeAPI {
        return CozeAPI.Builder()
            .auth(TokenAuth(Config.getInstance().getCozeAccessToken())) // 使用配置中的Token进行身份验证
            .baseURL(Config.getInstance().baseURL) // 设置基础URL
            .build() // 构建CozeAPI实例
    }

    /**
     * 提供初始化后的CozeAPI实例
     *
     * @return CozeAPI实例
     */
    fun getCozeAPI(): CozeAPI = cozeAPI

    /**
     * 伴生对象，用于实现CozeAPIManager的单例模式
     */
    companion object {
        @Volatile
        private var instance: CozeAPIManager? = null

        /**
         * 获取CozeAPIManager单例实例的方法
         * 如果实例不存在，则同步创建新实例
         *
         * @return CozeAPIManager单例实例
         */
        fun getInstance(): CozeAPIManager =
            instance ?: synchronized(this) {
                instance ?: CozeAPIManager().also { instance = it }
            }
    }
}
