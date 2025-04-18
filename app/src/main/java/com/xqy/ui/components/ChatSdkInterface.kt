package com.xqy.ui.components

import android.annotation.SuppressLint
import android.content.Context
import android.webkit.WebView
import android.webkit.WebViewClient
import android.webkit.WebChromeClient
import android.webkit.WebResourceRequest
import android.webkit.WebResourceError
import android.util.Log
import android.view.ViewGroup
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.viewinterop.AndroidView
import com.xqy.rtc.config.Config

/**
 * ChatSdkInterface 提供了一个直接集成Coze聊天SDK的接口
 * 它使用WebView但进行了优化，避免了额外的HTML页面加载
 */
class ChatSdkInterface(private val context: Context) {
    
    private var webView: WebView? = null
    
    /**
     * JavaScript接口类，用于向WebView注入token，getToken方法由外部js调用，所以这边显示没有调用，但其实有用
     */
    private inner class WebAppInterface {
        @android.webkit.JavascriptInterface
        fun getToken(): String {
            // 在非阻塞上下文中获取token
            var token: String
            kotlinx.coroutines.runBlocking {
                token = Config.getInstance().getCozeAccessToken()
            }
            return token
        }
    }
    
    /**
     * 初始化并配置WebView以直接加载聊天SDK
     */
    @SuppressLint("SetJavaScriptEnabled")
    fun initWebView(): WebView {
        return WebView(context).apply {
            webView = this
            
            // 配置WebView设置
            settings.apply {
                javaScriptEnabled = true
                domStorageEnabled = true
                allowFileAccess = true
                allowContentAccess = true
            }
            
            // 添加JavaScript接口
            addJavascriptInterface(WebAppInterface(), "Android")
            
            // 设置WebViewClient以处理页面加载
            webViewClient = object : WebViewClient() {
                override fun onPageFinished(view: WebView?, url: String?) {
                    super.onPageFinished(view, url)
                    // 页面加载完成后的处理
                    Log.d("ChatSdkInterface", "Page loaded successfully: $url")
                }
                
                override fun onReceivedError(view: WebView?, request: WebResourceRequest?, error: WebResourceError?) {
                    super.onReceivedError(view, request, error)
                    Log.e("ChatSdkInterface", "WebView error: ${error?.description}")
                }
            }
            
            // 设置WebChromeClient以处理JavaScript对话框
            webChromeClient = WebChromeClient()
            
            // 设置布局参数
            layoutParams = ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT
            )
            
            // 直接加载assets目录下的HTML文件
            loadUrl("file:///android_asset/chatsdk/index.html")
        }
    }

}

/**
 * 提供一个Composable函数来在UI中使用ChatSdkInterface
 */
@Composable
fun ChatSdkView(modifier: Modifier = Modifier) {
    AndroidView(
        factory = { context ->
            ChatSdkInterface(context).initWebView()
        },
        modifier = modifier
    )
}