
fn main() {
   let s = String::from("Hello");

    gotree::initialize(&s);
    gotree::mutate(&s);

    println!("Hello, world!");
}
