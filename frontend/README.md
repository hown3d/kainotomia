# Yew Styles Trunk Template

Template to start an app using [Yew Styles framework](https://yewstyles.spielrs.tech/) with [Trunk](https://trunkrs.dev/) as a builder.

# How it works

1. Install [Trunk](https://trunkrs.dev/#install)
2. Clone this template
3. Execute inside of the template: `trunk serve`

# Usage

## Installation
If you don't already have it installed, it's time to install Rust: https://www.rust-lang.org/tools/install. The rest of this guide assumes a typical Rust installation which contains both rustup and Cargo.

To compile Rust to WASM, we need to have the wasm32-unknown-unknown target installed. If you don't already have it, install it with the following command:

```
rustup target add wasm32-unknown-unknown
```
Now that we have our basics covered, it's time to install the star of the show: Trunk. Simply run the following command to install it:

```
cargo install --locked trunk  wasm-bindgen-cli
```
That's it, we're done!

## Running
```
trunk serve
```
Rebuilds the app whenever a change is detected and runs a local server to host it.

There's also the trunk watch command which does the same thing but without hosting it.

## Release
```
trunk build --release
```

This builds the app in release mode similar to cargo build --release. You can also pass the --release flag to trunk serve if you need to get every last drop of performance.

Unless overwritten, the output will be located in the dist directory.

# Directoies and files

* `src/assets` : images, icons and other assets are here
* `src/styles` : here are the `Yew Styles` css files and `main.sass`. **Note**: Only use `main.sass` or create new one in order to write your own styles
* `src/pages` : each component page represent a route
* `src/app.rs` : it is where the router is set and the common elements are added for the whole application
* `src/main.rs` : here call the yew instance in order to run the app


## License

Yew Styles Trunk Template is MIT licensed. See [license](LICENSE.md)
