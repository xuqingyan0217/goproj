package com.xqy.rtc.utils

import android.content.Context
import android.os.Handler
import android.os.Looper
import android.widget.Toast
import androidx.appcompat.app.AlertDialog

/**
 * ToastUtil 是一个用于显示 Toast 消息和对话框的工具类。
 * 它提供了在UI线程中显示长时Toast、短时Toast和错误对话框的方法。
 */
object ToastUtil {
    // uiHandler 用于在UI线程中执行操作，确保Toast和对话框在UI线程中显示。
    private val uiHandler = Handler(Looper.getMainLooper())
    // toast 用于保存当前显示的Toast，以便在显示新的Toast之前取消之前的Toast。
    private var toast: Toast? = null
    // dialog 用于保存当前显示的对话框，以防止重复显示对话框。
    private var dialog: AlertDialog? = null

    /**
     * 显示错误对话框。
     *
     * @param context 上下文，用于构建对话框。
     * @param message 要显示的错误信息。
     */
    fun showAlert(context: Context, message: String) {
        uiHandler.post {
            // 如果对话框已经在显示，则不执行任何操作。
            if (dialog?.isShowing == true) {
                return@post
            }
            // 构建并显示新的对话框。
            dialog = AlertDialog.Builder(context)
                .setTitle("错误")
                .setMessage(message)
                .setPositiveButton("OK") { dialog, _ -> dialog.dismiss() }
                .show()
        }
    }

    /**
     * 显示长时Toast消息。
     *
     * @param context 上下文，用于构建Toast。
     * @param msg 要显示的消息。
     */
    fun showLongToast(context: Context, msg: String) {
        uiHandler.post {
            // 取消之前的Toast，避免重复显示。
            toast?.cancel()
            // 构建并显示新的长时Toast消息。
            toast = Toast.makeText(context, msg, Toast.LENGTH_LONG)
            toast?.show()
        }
    }

    /**
     * 显示短时Toast消息。
     *
     * @param context 上下文，用于构建Toast。
     * @param msg 要显示的消息。
     */
    fun showShortToast(context: Context, msg: String) {
        uiHandler.post {
            // 取消之前的Toast，避免重复显示。
            toast?.cancel()
            // 构建并显示新的短时Toast消息。
            toast = Toast.makeText(context, msg, Toast.LENGTH_SHORT)
            toast?.show()
        }
    }
}
