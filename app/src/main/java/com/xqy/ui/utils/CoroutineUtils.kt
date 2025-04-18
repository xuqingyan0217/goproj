package com.xqy.ui.utils

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

/**
 * 协程工具类，用于处理后台任务
 */
object CoroutineUtils {
    /**
     * 在IO线程执行耗时操作
     * @param block 要执行的操作
     * @return 操作结果
     */
    suspend fun <T> executeOnBackground(block: suspend CoroutineScope.() -> T): T {
        return withContext(Dispatchers.IO) {
            block()
        }
    }
}