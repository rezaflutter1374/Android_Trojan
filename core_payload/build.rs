fn main() {
    // This build script is used for Android-specific configurations
    // Only link Android libraries when building for Android targets
    let target_os = std::env::var("CARGO_CFG_TARGET_OS").unwrap_or_default();
    
    if target_os == "android" {
        println!("cargo:rustc-link-lib=dylib=log");
        println!("cargo:rustc-link-lib=dylib=android");
        println!("cargo:rustc-link-lib=dylib=c");
        println!("Building for Android - linking Android libraries");
    } else {
        println!("Building for {} - skipping Android library linking", target_os);
    }
    
    // Tell Cargo that if the given file changes, to rerun this build script.
    println!("cargo:rerun-if-changed=src/main.rs");
}