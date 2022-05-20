use yew::prelude::*;

pub struct About;

impl Component for About {
    type Message = ();
    type Properties = ();

    fn create(_ctx: &Context<Self>) -> Self {
        About {}
    }

    fn update(&mut self, _ctx: &Context<Self>, _: Self::Message) -> bool {
        false
    }

    fn changed(&mut self, _ctx: &Context<Self>) -> bool {
        false
    }

    fn view(&self, _ctx: &Context<Self>) -> Html {
        html! {
            <div>{"hello"}</div>
        }
    }
}
