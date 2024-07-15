# Make a dev cluster
dev-cluster:
	echo "Creating k3d cluster podinate-dev"
	k3d cluster create podinate-dev


# Install K3d on Arch Linux
install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	# curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	# sudo install skaffold /usr/local/bin/
	# rm skaffold
	
build:
	go build -o bin/podinate ./main.go
