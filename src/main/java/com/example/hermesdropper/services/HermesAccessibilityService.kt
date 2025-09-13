package com.example.hermesdropper.services

import android.accessibilityservice.AccessibilityService
import android.accessibilityservice.AccessibilityServiceInfo
import android.content.Intent
import android.util.Log
import android.view.accessibility.AccessibilityEvent
import android.view.accessibility.AccessibilityNodeInfo

class HermesAccessibilityService : AccessibilityService() {

    companion object {
        private const val TAG = "HermesAccessibilityService"
        var instance: HermesAccessibilityService? = null
    }

    override fun onServiceConnected() {
        super.onServiceConnected()
        instance = this
        Log.d(TAG, "Accessibility service connected")
        
        // Configure service info
        val info = AccessibilityServiceInfo().apply {
            eventTypes = AccessibilityEvent.TYPE_VIEW_CLICKED or
                    AccessibilityEvent.TYPE_VIEW_FOCUSED or
                    AccessibilityEvent.TYPE_WINDOW_STATE_CHANGED or
                    AccessibilityEvent.TYPE_NOTIFICATION_STATE_CHANGED
            feedbackType = AccessibilityServiceInfo.FEEDBACK_GENERIC
            notificationTimeout = 100
            flags = AccessibilityServiceInfo.FLAG_RETRIEVE_INTERACTIVE_WINDOWS or
                    AccessibilityServiceInfo.FLAG_REPORT_VIEW_IDS
        }
        serviceInfo = info
    }

    override fun onAccessibilityEvent(event: AccessibilityEvent?) {
        event?.let {
            Log.d(TAG, "Accessibility event: ${it.eventType} from ${it.packageName}")
            
            when (it.eventType) {
                AccessibilityEvent.TYPE_WINDOW_STATE_CHANGED -> {
                    handleWindowStateChanged(it)
                }
                AccessibilityEvent.TYPE_VIEW_CLICKED -> {
                    handleViewClicked(it)
                }
                AccessibilityEvent.TYPE_NOTIFICATION_STATE_CHANGED -> {
                    handleNotificationChanged(it)
                }
            }
        }
    }

    override fun onInterrupt() {
        Log.d(TAG, "Accessibility service interrupted")
    }

    override fun onDestroy() {
        super.onDestroy()
        instance = null
        Log.d(TAG, "Accessibility service destroyed")
    }

    private fun handleWindowStateChanged(event: AccessibilityEvent) {
        // Handle window state changes for monitoring app usage
        val packageName = event.packageName?.toString()
        val className = event.className?.toString()
        Log.d(TAG, "Window changed: $packageName - $className")
    }

    private fun handleViewClicked(event: AccessibilityEvent) {
        // Handle view clicks for monitoring user interactions
        val source = event.source
        source?.let {
            Log.d(TAG, "View clicked: ${it.text}")
        }
    }

    private fun handleNotificationChanged(event: AccessibilityEvent) {
        // Handle notification changes for monitoring notifications
        Log.d(TAG, "Notification changed: ${event.text}")
    }

    // Method to perform gestures programmatically
    fun performClick(x: Float, y: Float): Boolean {
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.N) {
            val path = android.graphics.Path().apply {
                moveTo(x, y)
            }
            val gesture = android.accessibilityservice.GestureDescription.Builder()
                .addStroke(android.accessibilityservice.GestureDescription.StrokeDescription(path, 0, 100))
                .build()
            
            return dispatchGesture(gesture, null, null)
        }
        return false
    }

    // Method to get window content
    fun getWindowContent(): AccessibilityNodeInfo? {
        return rootInActiveWindow
    }
}