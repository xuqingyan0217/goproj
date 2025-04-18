package com.xqy.nutrition.data.model

/**
 * 四元组数据类
 * 用于同时返回四个不同类型的值
 */
data class Quadruple<A, B, C, D>(
    val first: A,
    val second: B,
    val third: C,
    val fourth: D
)