package com.xqy.ui.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import com.xqy.ui.screens.HomeScreen
import com.xqy.ui.screens.ScanScreen
import com.xqy.ui.screens.RecommendScreen
import com.xqy.ui.screens.ProfileScreen
import com.xqy.ui.screens.RTCScreen

object NavRoutes {
    const val HOME = "home"
    const val SCAN = "scan"
    const val RECOMMEND = "recommend"
    const val PROFILE = "profile"
    const val RTC = "rtc"
}

@Composable
fun NavGraph(
    navController: NavHostController,
    startDestination: String = NavRoutes.HOME
) {
    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable(NavRoutes.HOME) {
            HomeScreen()
        }
        composable(NavRoutes.SCAN) {
            ScanScreen()
        }
        composable(NavRoutes.RECOMMEND) {
            RecommendScreen()
        }
        composable(NavRoutes.RTC) {
            RTCScreen()
        }
        composable(NavRoutes.PROFILE) {
            ProfileScreen()
        }
    }
}