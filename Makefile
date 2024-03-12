# Make a dev cluster
dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	# Install the cockroach operator and crds 

# Install K3d on Arch Linux
install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	# curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	# sudo install skaffold /usr/local/bin/
	# rm skaffold
	
# Show the logs for the API backend in the Kubernetes cluster while developing
dev-backend-logs:
	kubectl -n podinate logs -l app=podinate-controller -f

dev-backend-shell:
	kubectl -n podinate exec -it deployment/podinate-controller -- /bin/bash

dev-code-upload:
	./api-backend/scripts/initial-code-upload.sh
	kubycat ./kubycat.yaml

make dev-port-forward:
	kubectl -n podinate port-forward service/podinate-controller 3001:3000

# Get a shell on the API backend Postgres pod (for debugging)
postgres-shell:
	bash -c "kubectl -n podinate exec -it postgres-0 -- psql 'postgresql://postgres:$$(kubectl -n podinate get secret masterdb-secret -o jsonpath='{.data.superUserPassword}' | base64 --decode ; echo)@localhost/podinate'"

# After API spec change, rebuild the generate code
api-generate:
	bash api/generate.sh

salt-sync:
	ssh ubuntu@salt.podinate.com "rm -rf ~/salt/*"
	scp -r infrastructure/salt/* ubuntu@salt.podinate.com:~/salt/
	ssh ubuntu@salt.podinate.com "sudo cp -r ~/salt/* /srv/salt/"

salt-apply: salt-sync
	ssh ubuntu@salt.podinate.com "sudo salt '*' state.apply"