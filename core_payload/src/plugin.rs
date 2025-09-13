use libloading::{Library, Symbol};
use serde_json::Value;
use std::collections::HashMap;
use std::path::Path;

pub type PluginFunction = unsafe fn(Value) -> Value;

pub struct Plugin {
    library: Library,
    functions: HashMap<String, Symbol<PluginFunction>>,
}

pub struct PluginManager {
    plugins: HashMap<String, Plugin>,
    next_id: u32,
}

impl PluginManager {
    pub fn new() -> Self {
        PluginManager {
            plugins: HashMap::new(),
            next_id: 1,
        }
    }

    pub fn load_plugin<P: AsRef<Path>>(&mut self, path: P) -> Result<String, String> {
        // Generate a unique ID for this plugin
        let plugin_id = format!("plugin_{}", self.next_id);
        self.next_id += 1;

        // Load the dynamic library
        let library = unsafe { Library::new(path.as_ref()) }
            .map_err(|e| format!("Failed to load plugin library: {}", e))?;

        // Create a new plugin instance
        let plugin = Plugin {
            library,
            functions: HashMap::new(),
        };

        // Store the plugin
        self.plugins.insert(plugin_id.clone(), plugin);

        Ok(plugin_id)
    }

    pub fn execute_plugin(
        &self,
        plugin_id: &str,
        function_name: &str,
        args: &Value,
    ) -> Result<Value, String> {
        // Find the plugin
        let plugin = self
            .plugins
            .get(plugin_id)
            .ok_or_else(|| format!("Plugin not found: {}", plugin_id))?;

        // Check if we've already loaded this function
        let function = if let Some(func) = plugin.functions.get(function_name) {
            func
        } else {
            // Load the function from the library
            let func = unsafe {
                plugin
                    .library
                    .get::<PluginFunction>(function_name.as_bytes())
                    .map_err(|e| format!("Failed to load function '{}': {}", function_name, e))?
            };

            // Store the function for future use (this requires a mutable reference, which we don't have)
            // For simplicity, we'll just return the function without caching it
            func
        };

        // Call the function
        let result = unsafe { function(args.clone()) };

        Ok(result)
    }
}