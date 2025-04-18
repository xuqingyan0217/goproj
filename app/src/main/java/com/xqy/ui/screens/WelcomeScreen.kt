package com.xqy.ui.screens

import androidx.compose.animation.core.*
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.pager.HorizontalPager
import androidx.compose.foundation.pager.rememberPagerState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.scale
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.launch

/**
 * 欢迎界面组件
 * 显示应用介绍的多个页面，支持左右滑动，最后一页显示开始使用按钮
 * @param onGetStarted 点击开始使用按钮的回调
 */
@Composable
fun WelcomeScreen(onGetStarted: () -> Unit) {
    // 创建Pager状态
    val pagerState = rememberPagerState(pageCount = { 4 })
    val coroutineScope = rememberCoroutineScope()
    
    // 欢迎页面数据
    val pages = listOf(
        WelcomePage(
            title = "营养健康小助手",
            description = "记录饮食，分析营养，定制健康生活方式",
            iconResId = android.R.drawable.ic_menu_compass
        ),
        WelcomePage(
            title = "饮食记录",
            description = "轻松记录每日饮食，追踪卡路里和营养成分",
            iconResId = android.R.drawable.ic_menu_edit
        ),
        WelcomePage(
            title = "营养分析",
            description = "智能分析饮食结构，提供个性化营养建议",
            iconResId = android.R.drawable.ic_menu_info_details
        ),
        WelcomePage(
            title = "健康生活",
            description = "定制专属健康计划，陪伴您的健康生活",
            iconResId = android.R.drawable.ic_menu_myplaces
        )
    )
    
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
    ) {
        // 水平滑动页面
        HorizontalPager(
            state = pagerState,
            modifier = Modifier.fillMaxSize()
        ) { page ->
            WelcomePageContent(
                page = pages[page],
                isLastPage = page == pages.size - 1,
                onGetStarted = onGetStarted
            )
        }
        
        // 页面指示器
        Box(
            modifier = Modifier
                .align(Alignment.BottomCenter)
                .padding(bottom = 32.dp)
        ) {
            Row(
                horizontalArrangement = Arrangement.Center,
                modifier = Modifier.padding(16.dp)
            ) {
                repeat(pages.size) { iteration ->
                    val color = if (pagerState.currentPage == iteration) {
                        MaterialTheme.colorScheme.primary
                    } else {
                        MaterialTheme.colorScheme.onSurface.copy(alpha = 0.2f)
                    }
                    Box(
                        modifier = Modifier
                            .padding(2.dp)
                            .size(8.dp)
                            .background(color = color, shape = CircleShape)
                    )
                }
            }
        }
        
        // 如果不是最后一页，显示跳过按钮在右上角
        if (pagerState.currentPage < pages.size - 1) {
            TextButton(
                onClick = { coroutineScope.launch { pagerState.animateScrollToPage(pages.size - 1) } },
                modifier = Modifier
                    .align(Alignment.TopEnd)
                    .padding(top = 32.dp, end = 16.dp)
            ) {
                Text(
                    text = "跳过",
                    color = MaterialTheme.colorScheme.primary,
                    fontWeight = FontWeight.Medium
                )
            }
        }
    }
}

/**
 * 欢迎页面内容
 */
@Composable
fun WelcomePageContent(
    page: WelcomePage,
    isLastPage: Boolean,
    onGetStarted: () -> Unit
) {
    // 创建动画效果
    val infiniteTransition = rememberInfiniteTransition(label = "welcome_animation")
    val scale by infiniteTransition.animateFloat(
        initialValue = 0.95f,
        targetValue = 1.05f,
        animationSpec = infiniteRepeatable(
            animation = tween(2000, easing = LinearEasing),
            repeatMode = RepeatMode.Reverse
        ),
        label = "scale_animation"
    )
    
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
        modifier = Modifier
            .fillMaxSize()
            .padding(32.dp)
    ) {
        // 应用图标
        Box(
            modifier = Modifier
                .size(120.dp)
                .scale(scale),
            contentAlignment = Alignment.Center
        ) {
            // 使用圆形背景
            Box(
                modifier = Modifier
                    .size(120.dp)
                    .background(
                        color = MaterialTheme.colorScheme.primaryContainer,
                        shape = CircleShape
                    )
            )
            
            // 图标
            Icon(
                painter = painterResource(id = page.iconResId),
                contentDescription = "页面图标",
                tint = MaterialTheme.colorScheme.primary,
                modifier = Modifier.size(64.dp)
            )
        }
        
        Spacer(modifier = Modifier.height(32.dp))
        
        // 页面标题
        Text(
            text = page.title,
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        // 页面描述
        Text(
            text = page.description,
            fontSize = 16.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(horizontal = 24.dp),
            textAlign = TextAlign.Center
        )
        
        Spacer(modifier = Modifier.height(48.dp))
        
        // 只在最后一页显示开始按钮
        if (isLastPage) {
            Button(
                onClick = onGetStarted,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(56.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.primary
                ),
                shape = RoundedCornerShape(28.dp)
            ) {
                Text(
                    text = "开始吧",
                    fontSize = 18.sp,
                    fontWeight = FontWeight.Medium
                )
            }
        }
        
        Spacer(modifier = Modifier.height(24.dp))
        
        // 版本信息
        if (isLastPage) {
            Text(
                text = "Version 1.0",
                fontSize = 12.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant.copy(alpha = 0.7f)
            )
        }
    }
}

/**
 * 欢迎页面数据类
 */
data class WelcomePage(
    val title: String,
    val description: String,
    val iconResId: Int
)