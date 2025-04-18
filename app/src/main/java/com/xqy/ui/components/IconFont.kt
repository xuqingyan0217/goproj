package com.xqy.ui.components

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.Font
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.unit.TextUnit
import androidx.compose.ui.unit.sp
import com.xqy.nutrition.R

/**
 * 使用IconFont字体图标
 *
 * @param iconId 图标ID，用于指定具体的图标
 * @param modifier 可选的修饰符，用于自定义图标的外观和布局
 * @param size 可选的字体大小，默认为24.sp
 * @param tint 可选的颜色，默认为未指定颜色
 */
@Composable
fun IconFont(
    iconId: String,
    modifier: Modifier = Modifier,
    size: TextUnit = 24.sp,
    tint: Color = Color.Unspecified
) {
    // 创建FontFamily对象，使用特定的图标字体文件
    val fontFamily = FontFamily(
        Font(R.font.iconfont)
    )

    // 显示图标文本，应用图标字体、大小和颜色
    Text(
        text = iconId,
        fontFamily = fontFamily,
        fontSize = size,
        color = tint,
        modifier = modifier
    )
}
