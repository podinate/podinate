dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	# Install the cockroach operator and crds 


install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	sudo install skaffold /usr/local/bin/
	rm skaffold
	
dev-code-api:
	./api-backend/scripts/initial-code-upload.sh

postgres-shell:
	bash -c "kubectl -n api exec -it masterdb-1-0 -- psql 'postgresql://postgres:$$(kubectl -n api get secret masterdb-secret -o jsonpath='{.data.superUserPassword}' | base64 --decode ; echo)@localhost/podinate'"