# CLI 
This is the official CLI for Podinate. It is written in Go, currently no way to distribute it is set up. 

## Cobra / Viper
The CLI is built on the Cobra framework / package, which handles all the routing, and also generates help and man pages for us pretty effectively. To see the help, run `podinate help`. 

## TUI
The CLI uses a TUI where possible when called in an interactive session. The TUI is built on the Bubble Tea package, using mostly Charms from the package written by the same authors, and a couple of third party charms where applicable. 

## Documentation
The CLI documents itself through the Cobra package. When you're creating a new command, make sure to add both long and short documentation for how to use your command. 

## Commands
These are the currently implemented commands. In Cobra, it's easy to set up shortened aliases to make commands more terse, but we should always refer to the full commands in documentation etc. 

### get / list
Get a list of something. So far I've only implemented `get pods` and `get projects`.

### launch
Launch launches a wizard which asks what you want to launch, such as for example an existing app, such as WordPress, or your own app, which would be the Dockerfile in the current directory. For MVP we are concentrating on launching existing apps as they're easy to test, then we will add launching custom apps, then we will add static sites. 

### tofu
The `tofu` sub-command passes through to OpenTofu or Terraform, depending which is installed. Any arguments provided are passed straight through to the child. For example you can run `podinate tf apply`. 