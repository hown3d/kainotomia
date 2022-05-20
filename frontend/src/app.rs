use crate::pages::{About, Home};
use yew::prelude::*;
use yew_router::prelude::*;


pub struct App {
    navbar_items: Vec<bool>,
}

#[derive(Routable, PartialEq, Clone, Debug, Copy)]
pub enum Route {
    #[at("/")]
    Home,
    #[at("/about")]
    About,
    #[not_found]
    #[at("/404")]
    NotFound,
}

fn switch(routes: &Route) -> Html {
    match routes {
        Route::Home => html! { <Home />},
        Route::About => html! { <About />},
        Route::NotFound => html! { <h1>{ "404" }</h1> },
    }
}

pub enum Msg {
    ChangeNavbarItem(usize),
}

impl Component for App {
    type Message = Msg;
    type Properties = ();

    fn create(_ctx: &Context<Self>) -> Self {
        App {
            navbar_items: vec![true, false],
        }
    }

    fn update(&mut self, _ctx: &Context<Self>, msg: Self::Message) -> bool {
        match msg {
            Msg::ChangeNavbarItem(index) => {
                for (i, _) in self.navbar_items.clone().into_iter().enumerate() {
                    self.navbar_items[i] = false;
                }

                self.navbar_items[index] = true;
            }
        }
        true
    }

    fn changed(&mut self, _ctx: &Context<Self>) -> bool {
        false
    }

    fn view(&self, _ctx: &Context<Self>) -> Html {
        html! {
        <BrowserRouter>
        <main>
            <Switch<Route> render={Switch::render(switch)} />
        </main>
        </BrowserRouter>
        }
    }

    fn rendered(&mut self, ctx: &Context<Self>, first_render: bool) {
        todo!()
    }

    fn destroy(&mut self, ctx: &Context<Self>) {
        todo!()
    }
}
