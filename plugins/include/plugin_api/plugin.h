// Plugin API header for C++ plugins
#ifndef PLUGIN_API_H
#define PLUGIN_API_H

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Execute the plugin action with the given payload.
 * 
 * @param payload A JSON string containing parameters for the action
 * @return A dynamically allocated string containing the result (must be freed by the caller)
 *         The string should be a JSON-encoded result
 */
char* execute_action(const char* payload);

#ifdef __cplusplus
}
#endif

#endif // PLUGIN_API_H