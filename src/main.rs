mod error;
mod ffi;

pub use error::*;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Tree {
        pub root_cid: String,
}

pub fn initialize(argz: &str) -> Result<Tree>{
    let result: Result<Tree, String> = unsafe {
        ffi::call(ffi::bindings::Initialize, "Initialize", argz)?
    };

    result.map_err(Error::Fatal)

}

pub fn mutate(argz: &str){
    let result: Result<(), String> = unsafe { 
         ffi::call_raw_json(ffi::bindings::Mutate, "Mutate", argz).unwrap()
     };
    if let Err(err) = result {
        panic!("Failed to log in to twitter: {}", err);
    }
}

    fn main() {
	let s = r#"unused-arg"#;
match initialize(s){
     Ok(sk) => {    
         println!("prollytree root cid: {0}", sk.root_cid);
         mutate(&s);
	},
Err(e) => {
		 println!("{e}");
	},

     }
     
}
