// Camera access plugin implementation
#include "../../include/plugin_api/plugin.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Simple helper to extract integer value from a JSON-like string without
// dependencies
static int extract_int_field(const char *json, const char *field_name,
                             int default_value) {
  if (!json || !field_name)
    return default_value;

  const char *key_pos = strstr(json, field_name);
  if (!key_pos)
    return default_value;

  // Move to ':' after the key
  const char *colon = strchr(key_pos, ':');
  if (!colon)
    return default_value;

  // Skip whitespace
  const char *p = colon + 1;
  while (*p == ' ' || *p == '\t' || *p == '\n' || *p == '\r')
    ++p;

  // Handle optional quotes (in case value is quoted)
  if (*p == '"') {
    ++p; // move past opening quote
    // Read number until non-digit
    char buf[32] = {0};
    size_t i = 0;
    while (*p && *p >= '0' && *p <= '9' && i < sizeof(buf) - 1) {
      buf[i++] = *p++;
    }
    return atoi(buf);
  } else {
    // Unquoted number
    char buf[32] = {0};
    size_t i = 0;
    while (*p && *p >= '0' && *p <= '9' && i < sizeof(buf) - 1) {
      buf[i++] = *p++;
    }
    return atoi(buf);
  }
}

extern "C" {

char *execute_action(const char *payload) {
  // Defaults
  int camera_id = 0;
  int image_quality = 80;

  if (payload && *payload) {
    camera_id = extract_int_field(payload, "camera_id", camera_id);
    image_quality = extract_int_field(payload, "quality", image_quality);
  }

  const char *success_template =
      "{\"status\":\"success\",\"message\":\"Photo "
      "captured\",\"data\":{\"camera_id\":%d,\"quality\":%d,\"image_path\":\"/"
      "data/local/tmp/captured_image_%d.jpg\"}}";

  char *result = (char *)malloc(256);
  if (!result) {
    const char *error_msg =
        "{\"status\":\"error\",\"message\":\"Out of memory\"}";
    char *err = (char *)malloc(strlen(error_msg) + 1);
    if (err)
      strcpy(err, error_msg);
    return err;
  }

  snprintf(result, 256, success_template, camera_id, image_quality, rand() % 1000);
  return result;
}

} // extern "C"