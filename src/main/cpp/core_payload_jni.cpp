#if __has_include(<jni.h>)
#include <jni.h>
#else
// Fallback definitions for IDE linting environments where the Android NDK
// headers are not available. These are minimal types/macros to satisfy parsing
// only.
typedef const struct JNINativeInterface_ *JNIEnv;
typedef void *jobject;
typedef int jint;
#ifndef JNIEXPORT
#define JNIEXPORT __attribute__((visibility("default")))
#endif
#ifndef JNICALL
#define JNICALL
#endif
#endif

#if __has_include(<android/log.h>)
#include <android/log.h>
#else
#ifndef ANDROID_LOG_INFO
#define ANDROID_LOG_INFO 4
#endif
#ifndef ANDROID_LOG_ERROR
#define ANDROID_LOG_ERROR 6
#endif
inline int __android_log_print(int, const char *, const char *, ...) {
  return 0;
}
#endif

#include <exception>

#define LOG_TAG "CorePayloadJNI"
#define LOGI(...) __android_log_print(ANDROID_LOG_INFO, LOG_TAG, __VA_ARGS__)
#define LOGE(...) __android_log_print(ANDROID_LOG_ERROR, LOG_TAG, __VA_ARGS__)

// Function prototype for the Rust entry point
extern "C" {
int rust_main();
}

extern "C" JNIEXPORT jint JNICALL
Java_com_example_hermesdropper_MainActivity_startPayload(JNIEnv *env,
                                                         jobject /* this */) {
  LOGI("Starting core payload from JNI");

  try {
    // Call the Rust entry point
    int result = rust_main();
    LOGI("Rust payload returned: %d", result);
    return result;
  } catch (const std::exception &e) {
    LOGE("Exception in native code: %s", e.what());
    // return -1;
    return -1;
  } catch (...) {
    LOGE("Unknown exception in native code");
    return -2;
    // throw std::runtime_error("Unknown exception in native code");
  }
}