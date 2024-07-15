# Contribute to Podinate
Thank you for your interest in contributing to Podinate. Podinate is an alternative to `kubectl apply` that gives users precise control over what's happening in their Kubernetes cluster. 

If you're looking for something to contribute, check out our [Github Project](https://github.com/orgs/podinate/projects/1) and see what interests you.


## Repo Structure
Podinate uses the ubiquitous [Cobra](https://github.com/spf13/cobra) library to manage command line and flags, including documentation.

- */cmd* holds the command line interface and top level commands of Podinate. For now there's just `apply` and `version`. 
- */engine* holds the brains of Podinate - code that takes Kubernetes or PodFile definitions and gets the Kubernetes cluster to that state. You *must* check the readme in that directory to understand how that works. 
- */docs* contains a mkdocs documentation site hosted at [docs.podinate.com](https://docs.podinate.com)
- */kube_client* - contains some helper code for connecting to a Kubernetes cluster. 
- */scripts* - a couple of helper build scripts. 

When a user runs *podinate apply file.pf*, the rootCmd in */cmd* is called, which then delegates to the ApplyCmd in apply.go. That command calls engine.Parse and gets an *engine.Package* back, which the command then calls *package.Apply()* on.  

## Documentation
Documentation is stored in the */docs* directory, which is a [MkDocs](https://www.mkdocs.org/) project using the [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) framework. 

The docs are deployed to CloudFlare pages using the `docs/pages-build.sh` script. 

## Getting started with development

### Creating Development Environment
To get a complete working Podinate server, you need to get a local Kubernetes environment. For this Podinate recommends using K3d which you can install like so:
#### Arch:
```bash
yay -S rancher-k3d-bin
```

#### curl | sudo bash
```bash
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```

### Clone the Git Repository
Of course, you have to check out the git repository first. 
```bash
# Using ssh key
git clone git@github.com:podinate/podinate.git
# Using https
git clone https://github.com/podinate/podinate.git
```

Then `cd` into the new directory:
```bash
cd podinate
```

### Create a New Development Cluster
Once K3d is installed, there is a Make script to create a development cluster and deploy the Podinate controller to it. Run that now. 
```bash
make dev-cluster
```
This will create a cluster in K3d that you can use for development and testing. 

## Get Started Developing
When I'm developing, I like to create a directory called "testapp" in the root of the project, then cd to it and create an alias to use `go run` for testing.

```bash
mkdir -p testapp
cd testapp 
alias podinate="go run ../"
```

Here's a basic PodFile to get you started: 
```hcl
podinate {
    package = "hello-world"
    namespace = "default"
}

pod "hello-world" {
    image = "ubuntu"
    tag = "latest" 
    command = [ "/bin/bash", "-c", "--" ]
    arguments = [ "while true; do echo 'Hello from Fedora!'; sleep 2; done;" ]
}
```

Copy the contents into `hello.pf` and apply with:
```bash
podinate apply hello.pf
```

If you're looking for something to start with, check out our [Github Project](https://github.com/orgs/podinate/projects/1) and see what interests you. 

The controller and cli communicate through client/server packages built off the OpenAPI spec in `api/`, if you make any changes to it, run `make api-generate` to rebuild it. 

## Build & Deployment System
Podinate is built by [GoReleaser](https://goreleaser.com/) whenever a new tag is created. Have a look at `.github/goreleaser.yml`. This pipeline will build the code into a bunch of different formats and upload them to a Github Release. It will also scan the output directory and add any `*.deb` files to an APT repo. This APT repo is finally uploaded to CloudFlare to make the release live for Debian based distros. 

There is also a separate `.github/workflows/build.yml` which uses Docker to build an image to be added to the Github Container Repo. I haven't found a way to manage tags nicely, so currently whatever is the latest on Github Main is yeehaw'ed to `ghcr.io/podinate/podinate:latest`. 

## Tests
Lol