package com.xqy.ui.screens
/*
* # 代码解释
该代码实现了一个音视频通话界面，主要功能包括：
1. **权限管理**：请求摄像头、录音等权限，未授权时提示用户。
2. **房间连接与断开**：通过CozeAPI创建房间，初始化RTC引擎并加入房间；支持断开连接。
3. **音视频控制**：提供开关视频、静音/取消静音的功能。
4. **消息处理**：接收并解析房间内的消息，动态更新UI显示内容。
5. **打断功能**：发送特定消息中断当前对话。  

# 控制流图
```mermaid
flowchart TD
    A[开始] --> B{是否已授权}
    B -->|否| C[请求权限]
    B -->|是| D[初始化界面]
    C --> E{用户是否授予权限}
    E -->|否| F[提示授权失败]
    E -->|是| D
    D --> G[等待用户操作]
    G --> H{点击连接按钮}
    H -->|否| I{点击其他按钮}
    I -->|打开视频| J[启动视频捕获]
    I -->|关闭视频| K[停止视频捕获]
    I -->|静音| L[停止音频捕获]
    I -->|取消静音| M[启动音频捕获]
    I -->|打断| N[发送打断消息]
    H -->|是| O[检查是否已连接]
    O -->|否| P[创建房间并连接]
    P --> Q[初始化RTC引擎]
    Q --> R[加入房间]
    R --> S[更新UI]
    O -->|是| T[断开连接]
    T --> U[销毁RTC引擎]
    U --> V[重置状态]
    S --> W[监听消息]
    W --> X{是否为有效消息}
    X -->|是| Y[更新消息显示]
    X -->|否| Z[忽略消息]
```

* */
import android.Manifest
import android.content.pm.PackageManager
import android.util.Log
import android.view.TextureView
import android.widget.FrameLayout
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.viewinterop.AndroidView
import androidx.core.content.ContextCompat
import com.coze.openapi.client.audio.rooms.CreateRoomReq
import com.coze.openapi.client.audio.rooms.CreateRoomResp
import com.coze.openapi.client.chat.model.ChatEventType
import com.coze.openapi.client.connversations.message.model.Message
import com.coze.openapi.service.service.CozeAPI
import com.fasterxml.jackson.core.type.TypeReference
import com.fasterxml.jackson.databind.ObjectMapper
import com.ss.bytertc.engine.*
import com.ss.bytertc.engine.data.StreamIndex
import com.ss.bytertc.engine.handler.IRTCRoomEventHandler
import com.ss.bytertc.engine.handler.IRTCVideoEventHandler
import com.ss.bytertc.engine.type.*
import com.xqy.rtc.config.Config
import com.xqy.rtc.manager.CozeAPIManager
import com.xqy.rtc.utils.ToastUtil
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import android.os.Handler
import android.os.Looper

@Composable
fun RTCScreen() {
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    val mapper = remember { ObjectMapper() }

    var rtcVideo by remember { mutableStateOf<RTCVideo?>(null) }
    var rtcRoom by remember { mutableStateOf<RTCRoom?>(null) }
    var roomInfo by remember { mutableStateOf<CreateRoomResp?>(null) }
    val cozeCli = remember { CozeAPIManager.getInstance().getCozeAPI() }

    var isVideoEnabled by remember { mutableStateOf(false) }
    var isAudioEnabled by remember { mutableStateOf(true) }
    var isConnected by remember { mutableStateOf(false) }
    var messageText by remember { mutableStateOf("") }
    var roomId by remember { mutableStateOf("") }
    var hasPermissions by remember { mutableStateOf(false) }
    
    // 创建一个可组合函数内的更新消息文本的函数
    val updateMessageText: (String) -> Unit = { text ->
        messageText += text
    }

    val permissionLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.RequestMultiplePermissions(),
        onResult = { permissions ->
            hasPermissions = permissions.all { it.value }
            if (!hasPermissions) {
                ToastUtil.showAlert(context, "请授予必要的权限以使用音视频功能")
            }
        }
    )

    LaunchedEffect(Unit) {
        val permissions = arrayOf(
            Manifest.permission.CAMERA,
            Manifest.permission.RECORD_AUDIO,
            Manifest.permission.INTERNET,
            Manifest.permission.MODIFY_AUDIO_SETTINGS
        )
        val permissionsToRequest = permissions.filter { permission ->
            ContextCompat.checkSelfPermission(context, permission) != PackageManager.PERMISSION_GRANTED
        }.toTypedArray()

        if (permissionsToRequest.isNotEmpty()) {
            permissionLauncher.launch(permissionsToRequest)
        } else {
            hasPermissions = true
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // 本地视频预览
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(300.dp)
                .background(Color.Black)
        ) {
            AndroidView(
                factory = { ctx ->
                    FrameLayout(ctx).apply {
                        id = android.R.id.content
                    }
                },
                modifier = Modifier.fillMaxSize(),
                update = { view ->
                    if (isConnected && rtcVideo != null && isVideoEnabled) {
                        view.removeAllViews()
                        val localTextureView = TextureView(context)
                        view.addView(localTextureView)
                        VideoCanvas().apply {
                            renderView = localTextureView
                            renderMode = VideoCanvas.RENDER_MODE_HIDDEN
                            rtcVideo?.setLocalVideoCanvas(StreamIndex.STREAM_INDEX_MAIN, this)
                        }
                    } else {
                        view.removeAllViews()
                    }
                }
            )
        }

        Spacer(modifier = Modifier.height(16.dp))

        // 房间ID显示
        Text(
            text = "房间ID: $roomId",
            style = MaterialTheme.typography.bodyLarge,
            modifier = Modifier.fillMaxWidth(),
            textAlign = TextAlign.Center
        )

        Spacer(modifier = Modifier.height(16.dp))

        // 控制按钮
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            Button(
                onClick = {
                    if (!hasPermissions) {
                        ToastUtil.showAlert(context, "请先授予必要的权限")
                        return@Button
                    }
                    if (!isConnected) {
                        scope.launch(Dispatchers.IO) {
                            doConnect(context, cozeCli) { newRoomInfo ->
                                roomInfo = newRoomInfo
                                roomId = newRoomInfo.roomID
                                createRTCEngine(context, newRoomInfo) { newRtcVideo ->
                                    rtcVideo = newRtcVideo
                                    if (isAudioEnabled) startVoice(newRtcVideo)
                                    createAndJoinRoom(newRtcVideo, newRoomInfo, updateMessageText) { newRtcRoom ->
                                        rtcRoom = newRtcRoom
                                        isConnected = true
                                    }
                                }
                            }
                        }
                    } else {
                        disconnect(rtcRoom, rtcVideo)
                        isConnected = false
                        roomId = ""
                    }
                },
                colors = ButtonDefaults.buttonColors(
                    containerColor = if (isConnected) Color.Red else MaterialTheme.colorScheme.primary
                )
            ) {
                Text(if (isConnected) "断开连接" else "连接")
            }

            Button(
                onClick = {
                    if (rtcVideo != null) {
                        if (isVideoEnabled) {
                            rtcVideo?.stopVideoCapture()
                        } else {
                            rtcVideo?.startVideoCapture()
                        }
                        isVideoEnabled = !isVideoEnabled
                    } else {
                        ToastUtil.showAlert(context, "请先连接")
                    }
                }
            ) {
                Text(if (isVideoEnabled) "关闭视频" else "打开视频")
            }

            Button(
                onClick = {
                    if (rtcVideo != null) {
                        if (isAudioEnabled) {
                            rtcVideo?.stopAudioCapture()
                        } else {
                            rtcVideo?.startAudioCapture()
                        }
                        isAudioEnabled = !isAudioEnabled
                    } else {
                        ToastUtil.showAlert(context, "请先连接")
                    }
                }
            ) {
                Text(if (isAudioEnabled) "静音" else "打开声音")
            }

            Button(
                onClick = {
                    if (rtcRoom == null) {
                        ToastUtil.showAlert(context, "请先连接")
                        return@Button
                    }
                    try {
                        val data = mapOf(
                            "id" to "event_1",
                            "event_type" to "conversation.chat.cancel",
                            "data" to "{}"
                        )
                        rtcRoom?.sendUserMessage(
                            roomInfo?.uid,
                            mapper.writeValueAsString(data),
                            MessageConfig.RELIABLE_ORDERED
                        )
                        ToastUtil.showShortToast(context, "打断成功")
                    } catch (e: Exception) {
                        ToastUtil.showShortToast(context, "打断失败")
                    }
                }
            ) {
                Text("打断")
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        // 消息显示区域
        Text(
            text = messageText,
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f)
                .background(MaterialTheme.colorScheme.surfaceVariant)
                .padding(16.dp),
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant
        )
    }
}



private fun doConnect(
    context: android.content.Context,
    cozeCli: CozeAPI,
    onSuccess: (CreateRoomResp) -> Unit
) {
    try {
        val req = CreateRoomReq.builder()
            .botID(Config.getInstance().botID)
            .voiceID(Config.getInstance().voiceID)
            .build()
        
        val roomInfo = cozeCli.audio().rooms().create(req)
        onSuccess(roomInfo)
    } catch (e: Exception) {
        ToastUtil.showAlert(context, "连接失败: ${e.message}")
    }
}

private fun disconnect(rtcRoom: RTCRoom?, rtcVideo: RTCVideo?) {
    rtcRoom?.apply {
        leaveRoom()
        destroy()
    }
    rtcVideo?.apply {
        stopAudioCapture()
        stopVideoCapture()
        RTCVideo.destroyRTCVideo()
    }
}

private fun startVoice(rtcVideo: RTCVideo?) {
    rtcVideo?.startAudioCapture()
}

private fun createRTCEngine(
    context: android.content.Context,
    roomInfo: CreateRoomResp,
    onSuccess: (RTCVideo) -> Unit
) {
    try {
        RTCVideo.destroyRTCVideo()
        Thread.sleep(100)

        val rtcVideo = RTCVideo.createRTCVideo(
            context.applicationContext,
            roomInfo.appID,
            object : IRTCVideoEventHandler() {
                override fun onWarning(warn: Int) {
                    Log.w("RTCScreen", "RTCVideo warning: $warn")
                }

                override fun onError(err: Int) {
                    Log.e("RTCScreen", "RTCVideo error: $err")
                }
            },
            null,
            null
        )

        onSuccess(rtcVideo)
    } catch (e: Exception) {
        Log.e("RTCScreen", "创建引擎失败", e)
        ToastUtil.showAlert(context, "创建引擎失败: ${e.message}")
    }
}

// 定义UI线程Handler用于更新UI
private val uiHandler = Handler(Looper.getMainLooper())

private fun createAndJoinRoom(
    rtcVideo: RTCVideo?,
    roomInfo: CreateRoomResp,
    updateMessageText: (String) -> Unit,
    onSuccess: (RTCRoom) -> Unit
) {
    rtcVideo?.createRTCRoom(roomInfo.roomID)?.apply {
        setRTCRoomEventHandler(object : IRTCRoomEventHandler() {
            override fun onRoomStateChanged(roomId: String, uid: String, state: Int, extraInfo: String) {
                Log.w("RTCScreen", "roomId:$roomId, uid:$uid, state:$state, extraInfo:$extraInfo")
            }

            override fun onUserMessageReceived(uid: String, message: String) {
                try {
                    Log.d("RTCScreen", "收到消息: $message")
                    
                    // 解析顶层JSON
                    val messageMap = ObjectMapper().readValue(
                        message,
                        object : TypeReference<Map<String, Any>>() {}
                    )
                    
                    // 检查事件类型
                    val eventType = messageMap["event_type"]?.toString()
                    Log.d("RTCScreen", "事件类型: $eventType")
                    
                    if (eventType == ChatEventType.CONVERSATION_MESSAGE_DELTA.value) {
                        // 获取data字段
                        val data = messageMap["data"]
                        if (data != null) {
                            try {
                                // 直接将data对象转换为Message对象
                                val msgObj = ObjectMapper().convertValue(data, Message::class.java)
                                
                                // 更新UI显示消息内容
                                val content = msgObj.content
                                if (content != null && content.isNotEmpty()) {
                                    // 在UI线程中更新消息文本，将新消息追加到现有消息后面
                                    uiHandler.post {
                                        updateMessageText(content)
                                    }
                                    Log.d("RTCScreen", "更新消息内容: $content")
                                }
                            } catch (e: Exception) {
                                Log.e("RTCScreen", "解析Message对象失败: ${e.message}")
                                // 尝试直接显示data内容
                                uiHandler.post {
                                    updateMessageText("收到消息: $message")
                                }
                            }
                        }
                    } else {
                        // 处理其他类型的消息
                        Log.d("RTCScreen", "收到其他类型消息: $eventType")
                    }
                } catch (e: Exception) {
                    Log.e("RTCScreen", "解析消息失败: ${e.message}", e)
                    // 显示原始消息
                    uiHandler.post {
                        updateMessageText("收到原始消息: $message")
                    }
                }
            }
        })

        val userInfo = UserInfo(roomInfo.uid, "")
        val roomConfig = RTCRoomConfig(
            ChannelProfile.CHANNEL_PROFILE_CHAT_ROOM,
            true, true, true
        )

        joinRoom(roomInfo.token, userInfo, roomConfig)
        onSuccess(this)
    }
}