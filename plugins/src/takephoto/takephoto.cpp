// TakePhoto plugin implementation
#include "../../include/plugin_api/plugin.h"
#include <stdlib.h>
#include <string.h>

#if defined(__ANDROID__)
#include <android/log.h>
#define LOG_TAG "TakePhotoPlugin"
#define LOGI(...) __android_log_print(ANDROID_LOG_INFO, LOG_TAG, __VA_ARGS__)
#else
#include <stdio.h>
#define LOGI(...)                                                              \
  do {                                                                         \
    fprintf(stdout, __VA_ARGS__);                                              \
    fprintf(stdout, "\n");                                                     \
  } while (0)
#endif

// Very simple parser for key=value pairs in payload (e.g., "quality=high")
static void extract_kv(const char *payload, const char *key, char *out,
                       size_t out_sz, const char *default_val) {
  if (!out || out_sz == 0)
    return;
  out[0] = '\0';
  if (!payload || !key) {
    if (default_val)
      strncpy(out, default_val, out_sz - 1), out[out_sz - 1] = '\0';
    return;
  }
  size_t key_len = strlen(key);
  const char *pos = strstr(payload, key);
  if (!pos) {
    if (default_val)
      strncpy(out, default_val, out_sz - 1), out[out_sz - 1] = '\0';
    return;
  }
  pos += key_len;
  if (*pos == '=')
    ++pos; // move past '=' if present
  // copy until delimiter or end
  size_t i = 0;
  while (pos[i] && pos[i] != ';' && pos[i] != ',' && pos[i] != '\n' &&
         pos[i] != ' ' && i < out_sz - 1) {
    out[i] = pos[i];
    ++i;
  }
  out[i] = '\0';
  if (i == 0 && default_val)
    strncpy(out, default_val, out_sz - 1), out[out_sz - 1] = '\0';
}

extern "C" {

char *execute_action(const char *payload) {
  // Log input
  LOGI("TakePhoto plugin executed with payload: %s",
       payload ? payload : "<null>");

  // Defaults and parse payload
  char quality[32];
  extract_kv(payload, "quality", quality, sizeof(quality), "medium");

  // NOTE: Real implementation using NDK Camera2 would go here.
  // For now, return a mock success JSON with a fake path.
  const char *image_path = "/data/local/tmp/photo.jpg";

  const char *template_json =
      "{\"status\":\"success\",\"message\":\"Photo "
      "taken\",\"data\":{\"image_path\":\"%s\",\"quality\":\"%s\"}}";

  // Compute size and allocate
  size_t needed =
      strlen(template_json) + strlen(image_path) + strlen(quality) + 1;
  char *result = (char *)malloc(needed);
  if (!result) {
    const char *err = "{\"status\":\"error\",\"message\":\"Out of memory\"}";
    char *res_err = (char *)malloc(strlen(err) + 1);
    if (res_err)
      strcpy(res_err, err);
    return res_err;
  }

#if defined(_MSC_VER)

  _snprintf(result, needed, template_json, image_path, quality);
#else
  snprintf(result, needed, template_json, image_path, quality);
  result[needed - 1] = '\0';
  //   LOGI("TackePhoto: %s", result);

#endif

  return result;
}

} // extern "C"