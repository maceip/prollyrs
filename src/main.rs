mod error;
mod ffi;

pub use error::*;

pub fn initialize(argz: &str){
    let result: Result<(), String> = unsafe {
        ffi::call(ffi::bindings::Initialize, "Initialize", argz).unwrap()
    };

    if let Err(err) = result {
        panic!("Failed to log in to twitter: {}", err);
    }
}

pub fn mutate(argz: &str){
    let result: Result<(), String> = unsafe { 
         ffi::call(ffi::bindings::Mutate, "Mutate", argz).unwrap()
     };

    if let Err(err) = result {
        panic!("Failed to log in to twitter: {}", err);
    }}

    fn main() {
        let s = String::from("Hello");
     
         initialize(&s);
         mutate(&s);
     
         println!("Hello, world!");
     }
     
