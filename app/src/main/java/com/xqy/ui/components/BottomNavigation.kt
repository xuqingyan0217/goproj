// 导入必要的Compose库和资源
package com.xqy.ui.components

import androidx.compose.foundation.layout.size
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.unit.dp
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import com.xqy.nutrition.R

/**
 * 底部导航栏组件
 *
 * @param selectedItem 当前选中的项索引，默认为0
 * @param onItemSelected 项点击时的回调函数，默认为空
 */
@Composable
fun BottomNavigation(
    selectedItem: Int = 0,
    onItemSelected: (Int) -> Unit = {}
) {
    // 创建一个导航栏，设置其阴影高度和容器颜色
    NavigationBar(
        tonalElevation = 0.dp,
        containerColor = MaterialTheme.colorScheme.surface
    ) {
        // 导航项列表，每个项包括图标、标签、选中状态和点击事件处理
        NavigationBarItem(
            icon = {
                // 使用自定义的IconFont组件显示图标
                IconFont(
                    iconId = stringResource(id = R.string.home),
                    modifier = Modifier.size(24.dp),
                    // 根据是否选中来设置图标的颜色
                    tint = if (selectedItem == 0) MaterialTheme.colorScheme.primary
                    else MaterialTheme.colorScheme.onSurfaceVariant
                )
            },
            label = { Text("首页") },
            selected = selectedItem == 0,
            onClick = { onItemSelected(0) },
            // 设置导航项的颜色配置
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = MaterialTheme.colorScheme.primary,
                selectedTextColor = MaterialTheme.colorScheme.primary,
                indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
            )
        )
        // 其他导航项配置类似，不再一一注释
        NavigationBarItem(
            icon = {
                IconFont(
                    iconId = stringResource(id = R.string.scan),
                    modifier = Modifier.size(24.dp),
                    tint = if (selectedItem == 1) MaterialTheme.colorScheme.primary
                    else MaterialTheme.colorScheme.onSurfaceVariant
                )
            },
            label = { Text("扫描") },
            selected = selectedItem == 1,
            onClick = { onItemSelected(1) },
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = MaterialTheme.colorScheme.primary,
                selectedTextColor = MaterialTheme.colorScheme.primary,
                indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
            )
        )
        NavigationBarItem(
            icon = {
                IconFont(
                    iconId = stringResource(id = R.string.recommend),
                    modifier = Modifier.size(24.dp),
                    tint = if (selectedItem == 2) MaterialTheme.colorScheme.primary
                    else MaterialTheme.colorScheme.onSurfaceVariant
                )
            },
            label = { Text("推荐") },
            selected = selectedItem == 2,
            onClick = { onItemSelected(2) },
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = MaterialTheme.colorScheme.primary,
                selectedTextColor = MaterialTheme.colorScheme.primary,
                indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
            )
        )
        NavigationBarItem(
            icon = {
                IconFont(
                    iconId = stringResource(id = R.string.call),
                    modifier = Modifier.size(24.dp),
                    tint = if (selectedItem == 3) MaterialTheme.colorScheme.primary
                    else MaterialTheme.colorScheme.onSurfaceVariant
                )
            },
            label = { Text("通话") },
            selected = selectedItem == 3,
            onClick = { onItemSelected(3) },
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = MaterialTheme.colorScheme.primary,
                selectedTextColor = MaterialTheme.colorScheme.primary,
                indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
            )
        )
        NavigationBarItem(
            icon = {
                IconFont(
                    iconId = stringResource(id = R.string.me),
                    modifier = Modifier.size(24.dp),
                    tint = if (selectedItem == 4) MaterialTheme.colorScheme.primary
                    else MaterialTheme.colorScheme.onSurfaceVariant
                )
            },
            label = { Text("我的") },
            selected = selectedItem == 4,
            onClick = { onItemSelected(4) },
            colors = NavigationBarItemDefaults.colors(
                selectedIconColor = MaterialTheme.colorScheme.primary,
                selectedTextColor = MaterialTheme.colorScheme.primary,
                indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
            )
        )
    }
}
