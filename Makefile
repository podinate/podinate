dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	# Install the cockroach operator and crds 
	kubectl apply -f https://raw.githubusercontent.com/cockroachdb/cockroach-operator/v2.10.0/install/crds.yaml
	kubectl apply -f https://raw.githubusercontent.com/cockroachdb/cockroach-operator/v2.10.0/install/operator.yaml

install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	sudo install skaffold /usr/local/bin/
	rm skaffold
	
dev-code-api:
	./api-backend/scripts/initial-code-upload.sh