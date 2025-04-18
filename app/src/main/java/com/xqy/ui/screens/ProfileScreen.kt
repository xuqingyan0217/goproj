package com.xqy.ui.screens

import android.annotation.SuppressLint
import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.*
import androidx.compose.animation.expandVertically
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.shrinkVertically
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.asImageBitmap
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Dialog
import com.xqy.nutrition.data.model.UserProfileManager
import com.xqy.ui.utils.FileUtils
import com.xqy.ui.utils.PermissionUtils
import kotlin.math.sin

@SuppressLint("SetJavaScriptEnabled")
@Composable
fun ProfileScreen() {
    // 临时存储拍照的URI
    var tempImageUri by remember { mutableStateOf<Uri?>(null) }
    var wavePhase by remember { mutableFloatStateOf(0f) }
    val infiniteTransition = rememberInfiniteTransition(label = "wave")
    val waveProgress by infiniteTransition.animateFloat(
        initialValue = 0f,
        targetValue = 1f,
        animationSpec = infiniteRepeatable(
            animation = tween(2000, easing = LinearEasing),
            repeatMode = RepeatMode.Restart
        ),
        label = "wave"
    )

    LaunchedEffect(waveProgress) {
        wavePhase = waveProgress * 2 * Math.PI.toFloat()
    }

    // 添加状态控制聊天区域的显示和隐藏
    var showChatArea by remember { mutableStateOf(false) }
    
    // 添加状态控制对话框的显示和隐藏
    var showUsernameDialog by remember { mutableStateOf(false) }
    var showWeightGoalDialog by remember { mutableStateOf(false) }
    var showImagePickerOptions by remember { mutableStateOf(false) }
    var showCaloriesDialog by remember { mutableStateOf(false) }
    var showProteinDialog by remember { mutableStateOf(false) }
    var showCarbsDialog by remember { mutableStateOf(false) }
    var showFatDialog by remember { mutableStateOf(false) }
    
    // 用户输入状态
    var usernameInput by remember { mutableStateOf("") }
    var weightGoalInput by remember { mutableStateOf("") }
    var caloriesInput by remember { mutableStateOf("") }
    var proteinInput by remember { mutableStateOf("") }
    var carbsInput by remember { mutableStateOf("") }
    var fatInput by remember { mutableStateOf("") }

    val waveColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f)
    val context = LocalContext.current
    
    // 初始化UserProfileManager
    val userProfileManager = remember { UserProfileManager.getInstance(context) }
    
    // 用户数据状态
    var username by remember { mutableStateOf(userProfileManager.getUsername()) }
    var weightGoal by remember { mutableStateOf(userProfileManager.getWeightGoal()) }
    var avatarBitmap by remember { mutableStateOf(userProfileManager.getAvatar()) }
    
    // 营养目标数据状态
    var calories by remember { mutableStateOf(userProfileManager.getCalories()) }
    var protein by remember { mutableStateOf(userProfileManager.getProtein()) }
    var carbs by remember { mutableStateOf(userProfileManager.getCarbs()) }
    var fat by remember { mutableStateOf(userProfileManager.getFat()) }
    
    // 相机权限请求
    val permissionLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { permissions ->
        val allGranted = permissions.values.all { it }
        if (allGranted) {
            showImagePickerOptions = true
        }
    }
    
    // 图片选择器
    val galleryLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.GetContent()
    ) { uri: Uri? ->
        uri?.let {
            userProfileManager.saveAvatar(it, context)
            avatarBitmap = userProfileManager.getAvatar()
        }
    }
    
    // 相机拍照
    val cameraLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.TakePicture()
    ) { success: Boolean ->
        if (success) {
            tempImageUri?.let {
                userProfileManager.saveAvatar(it, context)
                avatarBitmap = userProfileManager.getAvatar()
            }
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.verticalGradient(
                    colors = listOf(
                        waveColor,
                        MaterialTheme.colorScheme.surface
                    )
                )
            )
            .padding(16.dp)
    ) {
        // 用户信息卡片
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 12.dp)
                .drawBehind {
                    val path = Path()
                    val height = size.height
                    val width = size.width
                    val waveHeight = height * 0.1f

                    path.moveTo(0f, height * 0.7f)
                    for (x in 0..width.toInt() step 16) {
                        val y = sin(x * 0.03f + wavePhase) * waveHeight
                        path.lineTo(x.toFloat(), height * 0.7f + y)
                    }
                    path.lineTo(width, height)
                    path.lineTo(0f, height)
                    path.close()

                    drawPath(
                        path = path,
                        color = waveColor
                    )
                },
            elevation = CardDefaults.cardElevation(defaultElevation = 8.dp),
            shape = RoundedCornerShape(24.dp),
            colors = CardDefaults.cardColors(
                containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.9f)
            )
        ) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(20.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                // 头像 - 添加点击事件
                Box(
                    modifier = Modifier
                        .size(80.dp)
                        .clip(CircleShape)
                        .clickable {
                            if (PermissionUtils.hasRequiredPermissions(context)) {
                                showImagePickerOptions = true
                            } else {
                                permissionLauncher.launch(PermissionUtils.getRequiredPermissions())
                            }
                        }
                        .border(2.dp, MaterialTheme.colorScheme.primary, CircleShape),
                    contentAlignment = Alignment.Center
                ) {
                    if (avatarBitmap != null) {
                        Image(
                            bitmap = avatarBitmap!!.asImageBitmap(),
                            contentDescription = "用户头像",
                            modifier = Modifier.fillMaxSize()
                        )
                    } else {
                        Box(
                            modifier = Modifier
                                .fillMaxSize()
                                .background(
                                    brush = Brush.radialGradient(
                                        colors = listOf(
                                            MaterialTheme.colorScheme.primary,
                                            MaterialTheme.colorScheme.secondary
                                        )
                                    )
                                ),
                            contentAlignment = Alignment.Center
                        ) {
                            Text(
                                text = username.firstOrNull()?.toString() ?: "我",
                                color = MaterialTheme.colorScheme.onPrimary,
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold
                            )
                        }
                    }
                }
                
                Spacer(modifier = Modifier.width(16.dp))
                
                Column(modifier = Modifier.weight(1f)) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            text = username,
                            fontSize = 20.sp,
                            fontWeight = FontWeight.Bold,
                            modifier = Modifier.weight(1f)
                        )
                        IconButton(onClick = { 
                            usernameInput = username
                            showUsernameDialog = true 
                        }) {
                            Icon(
                                imageVector = Icons.Default.Edit,
                                contentDescription = "编辑用户名",
                                tint = MaterialTheme.colorScheme.primary
                            )
                        }
                    }
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            text = "目标: 减重 ${weightGoal}kg",
                            fontSize = 14.sp,
                            color = MaterialTheme.colorScheme.onSurfaceVariant,
                            modifier = Modifier.weight(1f)
                        )
                        IconButton(onClick = { 
                            weightGoalInput = weightGoal
                            showWeightGoalDialog = true 
                        }) {
                            Icon(
                                imageVector = Icons.Default.Edit,
                                contentDescription = "编辑目标",
                                tint = MaterialTheme.colorScheme.primary
                            )
                        }
                    }
                }
            }
        }
        
        // 内容区域：根据状态显示营养目标或聊天区域
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 8.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = if (!showChatArea) "营养目标" else "营养助手",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold,
                color = MaterialTheme.colorScheme.primary,
                modifier = Modifier.padding(vertical = 4.dp)
            )
            
            // 切换按钮
            Button(
                onClick = { showChatArea = !showChatArea },
                shape = RoundedCornerShape(16.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.primary
                )
            ) {
                Text(
                    text = if (!showChatArea) "打开助手" else "查看目标",
                    color = MaterialTheme.colorScheme.onPrimary
                )
            }
        }
        
        // 根据状态显示不同内容
        AnimatedVisibility(
            visible = !showChatArea,
            enter = fadeIn() + expandVertically(),
            exit = fadeOut() + shrinkVertically()
        ) {
            // 营养目标卡片
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp),
                shape = RoundedCornerShape(24.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.9f)
                )
            ) {
                Column(
                    modifier = Modifier.padding(20.dp)
                ) {
                    NutritionGoalItem(
                        label = "每日热量", 
                        value = calories, 
                        unit = "千卡",
                        onEditClick = {
                            caloriesInput = calories
                            showCaloriesDialog = true
                        }
                    )
                    HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))
                    NutritionGoalItem(
                        label = "蛋白质", 
                        value = protein, 
                        unit = "克",
                        onEditClick = {
                            proteinInput = protein
                            showProteinDialog = true
                        }
                    )
                    HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))
                    NutritionGoalItem(
                        label = "碳水化合物", 
                        value = carbs, 
                        unit = "克",
                        onEditClick = {
                            carbsInput = carbs
                            showCarbsDialog = true
                        }
                    )
                    HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))
                    NutritionGoalItem(
                        label = "脂肪", 
                        value = fat, 
                        unit = "克",
                        onEditClick = {
                            fatInput = fat
                            showFatDialog = true
                        }
                    )
                }
            }
        }
        
        // 聊天区域
        AnimatedVisibility(
            visible = showChatArea,
            enter = fadeIn() + expandVertically(),
            exit = fadeOut() + shrinkVertically()
        ) {
            // 聊天组件
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp),
                shape = RoundedCornerShape(24.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.9f)
                )
            ) {
                com.xqy.ui.components.ChatSdkView(
                    modifier = Modifier.fillMaxSize()
                )
            }
        }
    }
    
    // 用户名编辑对话框
    if (showUsernameDialog) {
        Dialog(onDismissRequest = { showUsernameDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改用户名",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = usernameInput,
                        onValueChange = { usernameInput = it },
                        label = { Text("用户名") },
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showUsernameDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (usernameInput.isNotBlank()) {
                                userProfileManager.setUsername(usernameInput)
                                username = usernameInput
                                showUsernameDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
    
    // 减重目标编辑对话框
    if (showWeightGoalDialog) {
        Dialog(onDismissRequest = { showWeightGoalDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改减重目标",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = weightGoalInput,
                        onValueChange = { weightGoalInput = it },
                        label = { Text("减重目标 (kg)") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showWeightGoalDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (weightGoalInput.isNotBlank()) {
                                userProfileManager.setWeightGoal(weightGoalInput)
                                weightGoal = weightGoalInput
                                showWeightGoalDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
    
    // 图片选择器选项对话框
    if (showImagePickerOptions) {
        Dialog(onDismissRequest = { showImagePickerOptions = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "选择头像来源",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    Button(
                        onClick = {
                            showImagePickerOptions = false
                            galleryLauncher.launch("image/*")
                        },
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(vertical = 8.dp)
                    ) {
                        Text("从相册选择")
                    }
                    
                    Button(
                        onClick = {
                            showImagePickerOptions = false
                            val (_, uri) = FileUtils.createImageFile(context)
                            tempImageUri = uri
                            cameraLauncher.launch(uri)
                        },
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(vertical = 8.dp)
                    ) {
                        Text("拍照")
                    }
                    
                    TextButton(
                        onClick = { showImagePickerOptions = false },
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(vertical = 8.dp)
                    ) {
                        Text("取消")
                    }
                }
            }
        }
    }
    
    // 每日热量编辑对话框
    if (showCaloriesDialog) {
        Dialog(onDismissRequest = { showCaloriesDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改每日热量目标",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = caloriesInput,
                        onValueChange = { caloriesInput = it },
                        label = { Text("每日热量 (千卡)") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showCaloriesDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (caloriesInput.isNotBlank()) {
                                userProfileManager.setCalories(caloriesInput)
                                calories = caloriesInput
                                showCaloriesDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
    
    // 蛋白质编辑对话框
    if (showProteinDialog) {
        Dialog(onDismissRequest = { showProteinDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改蛋白质目标",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = proteinInput,
                        onValueChange = { proteinInput = it },
                        label = { Text("蛋白质 (克)") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showProteinDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (proteinInput.isNotBlank()) {
                                userProfileManager.setProtein(proteinInput)
                                protein = proteinInput
                                showProteinDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
    
    // 碳水化合物编辑对话框
    if (showCarbsDialog) {
        Dialog(onDismissRequest = { showCarbsDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改碳水化合物目标",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = carbsInput,
                        onValueChange = { carbsInput = it },
                        label = { Text("碳水化合物 (克)") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showCarbsDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (carbsInput.isNotBlank()) {
                                userProfileManager.setCarbs(carbsInput)
                                carbs = carbsInput
                                showCarbsDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
    
    // 脂肪编辑对话框
    if (showFatDialog) {
        Dialog(onDismissRequest = { showFatDialog = false }) {
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(
                        text = "修改脂肪目标",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(bottom = 16.dp)
                    )
                    
                    OutlinedTextField(
                        value = fatInput,
                        onValueChange = { fatInput = it },
                        label = { Text("脂肪 (克)") },
                        singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 16.dp),
                        horizontalArrangement = Arrangement.End
                    ) {
                        TextButton(onClick = { showFatDialog = false }) {
                            Text("取消")
                        }
                        Spacer(modifier = Modifier.width(8.dp))
                        Button(onClick = {
                            if (fatInput.isNotBlank()) {
                                userProfileManager.setFat(fatInput)
                                fat = fatInput
                                showFatDialog = false
                            }
                        }) {
                            Text("确定")
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun NutritionGoalItem(label: String, value: String, unit: String, onEditClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            text = label,
            fontSize = 16.sp
        )
        Row(verticalAlignment = Alignment.CenterVertically) {
            Row(verticalAlignment = Alignment.Bottom) {
                Text(
                    text = value,
                    fontSize = 18.sp,
                    fontWeight = FontWeight.Bold
                )
                Spacer(modifier = Modifier.width(4.dp))
                Text(
                    text = unit,
                    fontSize = 14.sp,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
            IconButton(onClick = onEditClick) {
                Icon(
                    imageVector = Icons.Default.Edit,
                    contentDescription = "编辑$label",
                    tint = MaterialTheme.colorScheme.primary
                )
            }
        }
    }
}