dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	alias k=kubectl
	k cluster-info

install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	sudo install skaffold /usr/local/bin/
	rm skaffold
	