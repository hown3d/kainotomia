#![recursion_limit = "1024"]

extern crate yew;
extern crate yew_router;

mod app;
mod pages;
use app::App;

pub fn main() {
    yew::start_app::<App>();
}
