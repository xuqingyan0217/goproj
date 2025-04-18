package com.xqy.ui.screens

import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import kotlinx.coroutines.launch
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.xqy.ui.utils.FileUtils
import com.xqy.ui.utils.PermissionUtils
import androidx.compose.foundation.Image
import androidx.compose.ui.layout.ContentScale
import android.graphics.BitmapFactory
import androidx.compose.ui.graphics.asImageBitmap
import androidx.compose.runtime.remember
import com.xqy.nutrition.data.network.ApiService
import com.xqy.ui.utils.CoroutineUtils
import com.xqy.ui.utils.ImageUtils
import android.widget.Toast
import androidx.compose.ui.res.stringResource
import com.xqy.nutrition.R
import com.xqy.ui.components.IconFont
import java.util.Locale
import com.xqy.ui.components.NutritionItem

@Composable
fun ScanScreen() {
    val context = LocalContext.current
    var selectedImageUri by remember { mutableStateOf<Uri?>(null) }
    var photoUri by remember { mutableStateOf<Uri?>(null) }
    var isAnalyzing by remember { mutableStateOf(false) }
    var analyzedFood by remember { mutableStateOf<String?>(null) }
    
    // 营养数据状态
    var calories by remember { mutableStateOf("0") }
    var protein by remember { mutableStateOf("0g") }
    var fat by remember { mutableStateOf("0g") }
    var carbs by remember { mutableStateOf("0g") }
    var isNutritionDataAvailable by remember { mutableStateOf(false) }

    // 图片选择器启动器
    val imagePickerLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.GetContent()
    ) { uri: Uri? ->
        uri?.let { 
            selectedImageUri = it
            // 开始分析过程
            isAnalyzing = true
            // 重置之前的分析结果
            analyzedFood = null
            isNutritionDataAvailable = false
            
            // 调用火山引擎API识别食物
            kotlinx.coroutines.MainScope().launch {
                try {
                    // 将图片转为Base64
                    val base64Image = CoroutineUtils.executeOnBackground {
                        ImageUtils.uriToBase64(context, uri)
                    }
                    
                    if (base64Image != null) {
                        // 调用API识别食物
                        val result = CoroutineUtils.executeOnBackground {
                            ApiService.recognizeFood(base64Image)
                        }
                        
                        result.onSuccess { nutrition ->
                            // 记录日志，跟踪营养数据
                            android.util.Log.d("ScanScreen", "收到API响应 - 食物: ${nutrition.name}, 卡路里: ${nutrition.calories}, 蛋白质: ${nutrition.protein}, 脂肪: ${nutrition.fat}")
                            
                            // 更新UI显示识别结果
                            analyzedFood = nutrition.name
                            calories = nutrition.calories.toString()
                            protein = "${String.format(Locale.CHINA,"%.1f", nutrition.protein)}g"
                            fat = "${String.format(Locale.CHINA,"%.1f", nutrition.fat)}g"
                            // 设置默认碳水值，实际应用中可能需要从API获取
                            val carbsValue = nutrition.calories * 0.15f / 4f // 假设碳水占总热量的15%，1g碳水产生4卡路里
                            carbs = "${String.format(Locale.CHINA,"%.1f", carbsValue)}g"
                            isNutritionDataAvailable = true
                            
                            // 添加到食物记录管理器前记录日志
                            android.util.Log.d("ScanScreen", "准备添加到FoodRecordManager - 蛋白质: ${nutrition.protein}, 脂肪: ${nutrition.fat}, 碳水: $carbsValue")
                            
                            // 添加到食物记录管理器
                            com.xqy.nutrition.data.model.FoodRecordManager.addRecord(
                                name = nutrition.name,
                                calories = nutrition.calories,
                                protein = nutrition.protein,
                                fat = nutrition.fat,
                                carbs = carbsValue
                            )
                        }.onFailure { error ->
                            // 显示错误信息
                            Toast.makeText(context, "识别失败: ${error.message}", Toast.LENGTH_SHORT).show()
                        }
                    } else {
                        Toast.makeText(context, "图片处理失败", Toast.LENGTH_SHORT).show()
                    }
                } catch (e: Exception) {
                    Toast.makeText(context, "发生错误: ${e.message}", Toast.LENGTH_SHORT).show()
                } finally {
                    isAnalyzing = false
                }
            }
        }
    }

    // 权限请求启动器和回调函数
    var permissionRequestCallback by remember { mutableStateOf<(() -> Unit)?>(null) }
    
    // 使用remember确保权限请求启动器在重组时保持一致
    val permissionLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.RequestMultiplePermissions()
    ) { permissions ->
        // 记录日志，跟踪权限请求结果
        android.util.Log.d("ScanScreen", "权限请求结果: $permissions")
        
        val allGranted = permissions.values.all { it }
        if (allGranted) {
            // 权限已获取，执行回调函数
            android.util.Log.d("ScanScreen", "所有权限已获取，执行回调")
            permissionRequestCallback?.invoke()
            permissionRequestCallback = null
        } else {
            // 权限被拒绝，显示提示
            Toast.makeText(context, "需要相机和存储权限才能使用此功能", Toast.LENGTH_SHORT).show()
            // 清除回调，避免在权限被拒绝的情况下执行操作
            permissionRequestCallback = null
        }
    }

    // 相机启动器
    val cameraLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.TakePicture()
    ) { success ->
        if (success) {
            // 处理拍照结果
            photoUri?.let { uri ->
                selectedImageUri = uri
                // 开始分析过程
                isAnalyzing = true
                // 重置之前的分析结果
                analyzedFood = null
                isNutritionDataAvailable = false
                
                // 调用火山引擎API识别食物
                kotlinx.coroutines.MainScope().launch {
                    try {
                        // 将图片转为Base64
                        val base64Image = CoroutineUtils.executeOnBackground {
                            ImageUtils.uriToBase64(context, uri)
                        }
                        
                        if (base64Image != null) {
                            // 调用API识别食物
                            val result = CoroutineUtils.executeOnBackground {
                                ApiService.recognizeFood(base64Image)
                            }
                            
                            result.onSuccess { nutrition ->
                                // 记录日志，跟踪营养数据
                                android.util.Log.d("ScanScreen", "收到API响应 - 食物: ${nutrition.name}, 卡路里: ${nutrition.calories}, 蛋白质: ${nutrition.protein}, 脂肪: ${nutrition.fat}")
                                // 更新UI显示识别结果
                                analyzedFood = nutrition.name
                                calories = nutrition.calories.toString()
                                protein = "${String.format(Locale.CHINA,"%.1f", nutrition.protein)}g"
                                fat = "${String.format(Locale.CHINA,"%.1f", nutrition.fat)}g"
                                // 设置默认碳水值，实际应用中可能需要从API获取
                                val carbsValue = nutrition.calories * 0.15f / 4f // 假设碳水占总热量的15%，1g碳水产生4卡路里
                                carbs = "${String.format(Locale.CHINA,"%.1f", carbsValue)}g"
                                isNutritionDataAvailable = true
                                // 添加到食物记录管理器前记录日志
                                android.util.Log.d("ScanScreen", "准备添加到FoodRecordManager - 蛋白质: ${nutrition.protein}, 脂肪: ${nutrition.fat}, 碳水: $carbsValue")

                                // 添加到食物记录管理器
                                com.xqy.nutrition.data.model.FoodRecordManager.addRecord(
                                    name = nutrition.name,
                                    calories = nutrition.calories,
                                    protein = nutrition.protein,
                                    fat = nutrition.fat,
                                    carbs = carbsValue
                                )
                            }.onFailure { error ->
                                // 显示错误信息
                                Toast.makeText(context, "识别失败: ${error.message}", Toast.LENGTH_SHORT).show()
                            }
                        } else {
                            Toast.makeText(context, "图片处理失败", Toast.LENGTH_SHORT).show()
                        }
                    } catch (e: Exception) {
                        Toast.makeText(context, "发生错误: ${e.message}", Toast.LENGTH_SHORT).show()
                    } finally {
                        isAnalyzing = false
                    }
                }
            }
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Top
    ) {
        // 相机预览区域 - 调整大小和位置
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .weight(0.4f)
                .padding(horizontal = 16.dp, vertical = 8.dp)
                .background(
                    color = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.9f),
                    shape = MaterialTheme.shapes.large
                ),
            contentAlignment = Alignment.Center
        ) {
            // 相机预览占位或选中的图片
            if (selectedImageUri != null) {
                val bitmap = remember(selectedImageUri) {
                    selectedImageUri?.let {
                        context.contentResolver.openInputStream(it)?.use { stream ->
                            BitmapFactory.decodeStream(stream)
                        }?.asImageBitmap()
                    }
                }
                bitmap?.let { imageBitmap ->
                    Box(modifier = Modifier.fillMaxSize()) {
                        Image(
                            bitmap = imageBitmap,
                            contentDescription = "选中的图片",
                            modifier = Modifier.fillMaxSize(),
                            contentScale = ContentScale.Fit
                        )
                        
                        // 显示分析中状态
                        if (isAnalyzing) {
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .background(Color.Black.copy(alpha = 0.5f)),
                                contentAlignment = Alignment.Center
                            ) {
                                CircularProgressIndicator(color = MaterialTheme.colorScheme.primary)
                            }
                        }
                        
                        // 显示识别结果
                        analyzedFood?.let { food ->
                            Box(
                                modifier = Modifier
                                    .align(Alignment.BottomStart)
                                    .padding(16.dp)
                                    .background(
                                        color = MaterialTheme.colorScheme.primaryContainer,
                                        shape = MaterialTheme.shapes.small
                                    )
                                    .padding(8.dp)
                            ) {
                                Text(
                                    text = "识别结果: $food",
                                    color = MaterialTheme.colorScheme.onPrimaryContainer,
                                    fontWeight = FontWeight.Bold
                                )
                            }
                        }
                    }
                }
            } else {
                IconFont(
                    iconId = stringResource(id = R.string.focus),
                    modifier = Modifier
                        .size(64.dp)
                        .background(
                            color = MaterialTheme.colorScheme.primary.copy(alpha = 0.08f),
                            shape = MaterialTheme.shapes.large
                        )
                        .padding(12.dp),
                    size = (40.sp),
                    tint = MaterialTheme.colorScheme.primary,
                )
            }
        }

        // 营养分析结果区域
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 8.dp),
            elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
            colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
        ) {
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        text = "营养分析结果",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold
                    )
                    
                    // 分析状态指示器
                    if (isAnalyzing) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            CircularProgressIndicator(
                                modifier = Modifier.size(16.dp),
                                strokeWidth = 2.dp
                            )
                            Spacer(modifier = Modifier.width(8.dp))
                            Text(
                                text = "分析中...",
                                fontSize = 14.sp,
                                color = MaterialTheme.colorScheme.primary
                            )
                        }
                    }
                }
                
                Spacer(modifier = Modifier.height(8.dp))
                
                // 识别的食物名称
                analyzedFood?.let { food ->
                    Text(
                        text = "识别食物: $food",
                        fontSize = 16.sp,
                        fontWeight = FontWeight.Medium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        modifier = Modifier.padding(bottom = 8.dp)
                    )
                }
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    NutritionItem(calories, "卡路里", MaterialTheme.colorScheme.primary)
                    NutritionItem(protein, "蛋白质", MaterialTheme.colorScheme.secondary)
                    NutritionItem(fat, "脂肪", MaterialTheme.colorScheme.tertiary)
                    NutritionItem(carbs, "碳水", MaterialTheme.colorScheme.error)
                }
                
                // 如果没有数据，显示提示
                if (!isNutritionDataAvailable && !isAnalyzing && selectedImageUri == null) {
                    Text(
                        text = "拍照或选择图片以分析食物营养成分",
                        fontSize = 14.sp,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        modifier = Modifier.padding(top = 8.dp)
                    )
                }
            }
        }
        
        // 底部控制栏
        Surface(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 8.dp),
            color = MaterialTheme.colorScheme.surface,
            shadowElevation = 4.dp,
            shape = MaterialTheme.shapes.medium
        ) {
            Column(
                modifier = Modifier.padding(16.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = "将食物放入取景框",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Medium,
                    color = MaterialTheme.colorScheme.onSurface
                )
                Spacer(modifier = Modifier.height(16.dp))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    // 相册选择按钮
                    Button(
                        onClick = {
                            android.util.Log.d("ScanScreen", "相册选择按钮点击")
                            if (!PermissionUtils.hasRequiredPermissions(context)) {
                                android.util.Log.d("ScanScreen", "需要请求权限")
                                permissionRequestCallback = { 
                                    android.util.Log.d("ScanScreen", "权限获取成功，启动相册选择器")
                                    imagePickerLauncher.launch("image/*") 
                                }
                                // 一次性请求所有需要的权限
                                permissionLauncher.launch(PermissionUtils.getRequiredPermissions())
                            } else {
                                android.util.Log.d("ScanScreen", "已有权限，直接启动相册选择器")
                                imagePickerLauncher.launch("image/*")
                            }
                        },
                        modifier = Modifier
                            .weight(1f)
                            .height(56.dp),
                        colors = ButtonDefaults.buttonColors(
                            containerColor = MaterialTheme.colorScheme.primary
                        )
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            IconFont(
                                iconId = stringResource(id = R.string.photo),
                                modifier = Modifier.size(24.dp),
                                tint = MaterialTheme.colorScheme.onPrimary
                            )
                            Spacer(modifier = Modifier.width(8.dp))
                            Text("从相册选择")
                        }
                    }
                    Spacer(modifier = Modifier.width(8.dp))
                    // 拍照按钮
                    Button(
                        onClick = {
                            android.util.Log.d("ScanScreen", "拍照按钮点击")
                            // 定义拍照操作
                            val takePictureAction = {
                                // 重置之前的分析结果
                                analyzedFood = null
                                isNutritionDataAvailable = false
                                
                                android.util.Log.d("ScanScreen", "创建图片文件并启动相机")
                                val (_, uri) = FileUtils.createImageFile(context)
                                photoUri = uri
                                cameraLauncher.launch(uri)
                            }
                            
                            if (!PermissionUtils.hasRequiredPermissions(context)) {
                                android.util.Log.d("ScanScreen", "需要请求权限才能拍照")
                                permissionRequestCallback = takePictureAction
                                // 一次性请求所有需要的权限
                                permissionLauncher.launch(PermissionUtils.getRequiredPermissions())
                            } else {
                                android.util.Log.d("ScanScreen", "已有权限，直接拍照")
                                takePictureAction()
                            }
                        },
                        modifier = Modifier
                            .weight(1f)
                            .height(56.dp),
                        colors = ButtonDefaults.buttonColors(
                            containerColor = MaterialTheme.colorScheme.primary
                        )
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            IconFont(
                                iconId = stringResource(id = R.string.camera),
                                modifier = Modifier.size(24.dp),
                                tint = MaterialTheme.colorScheme.onPrimary
                            )
                            Spacer(modifier = Modifier.width(8.dp))
                            Text("拍照识别")
                        }
                    }
                }
            }
        }
        
    }
}
