package com.xqy

import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.navigation.compose.rememberNavController
import com.xqy.nutrition.data.model.FoodRecordManager
import com.xqy.ui.components.BottomNavigation
import com.xqy.ui.navigation.NavGraph
import com.xqy.ui.navigation.NavRoutes
import com.xqy.ui.screens.WelcomeScreen
import com.xqy.ui.theme.NutritionTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // 初始化FoodRecordManager
        FoodRecordManager.initialize(applicationContext)
        Log.d("MainActivity", "FoodRecordManager initialized")
        
        // 初始化Config
        com.xqy.rtc.config.Config.init(applicationContext)
        Log.d("MainActivity", "Config initialized")
        
        setContent {
            NutritionTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    // 添加状态管理，控制是否显示欢迎界面
                    var showWelcomeScreen by remember { mutableStateOf(true) }
                    
                    // 创建导航控制器
                    val navController = rememberNavController()
                    var selectedTab by remember { mutableIntStateOf(0) }
                    
                    // 使用AnimatedVisibility实现动画过渡
                    AnimatedVisibility(
                        visible = showWelcomeScreen,
                        enter = fadeIn(),
                        exit = fadeOut(animationSpec = tween(durationMillis = 500))
                    ) {
                        // 显示欢迎界面
                        WelcomeScreen(onGetStarted = { showWelcomeScreen = false })
                    }
                    
                    AnimatedVisibility(
                        visible = !showWelcomeScreen,
                        enter = slideInVertically(initialOffsetY = { it / 2 }) + 
                                fadeIn(initialAlpha = 0.3f),
                        exit = slideOutVertically() + fadeOut()
                    ) {
                        // 显示主应用界面
                        
                        Scaffold(
                            bottomBar = {
                                BottomNavigation(
                                    selectedItem = selectedTab,
                                    onItemSelected = { index ->
                                        selectedTab = index
                                        when (index) {
                                            0 -> navController.navigate(NavRoutes.HOME)
                                            1 -> navController.navigate(NavRoutes.SCAN)
                                            2 -> navController.navigate(NavRoutes.RECOMMEND)
                                            3 -> navController.navigate(NavRoutes.RTC)
                                            4 -> navController.navigate(NavRoutes.PROFILE)
                                        }
                                    }
                                )
                            }
                        ) { paddingValues ->
                            Box(modifier = Modifier.padding(paddingValues)) {
                                NavGraph(navController = navController)
                            }
                        }
                    }
                }
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        // 清理Config和EnvConfig资源
        com.xqy.rtc.config.Config.cleanup()
        com.xqy.nutrition.data.auth.EnvConfig.cleanup()
        Log.d("MainActivity", "Config and EnvConfig cleaned up")
    }
}