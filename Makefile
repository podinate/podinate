# Make a dev cluster
dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	# Install the cockroach operator and crds 

# Install K3d on Arch Linux
install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	sudo install skaffold /usr/local/bin/
	rm skaffold
	
# Show the logs for the API backend in the Kubernetes cluster while developing
dev-backend-logs:
	kubectl -n api logs -l app=api-backend -f

dev-code-upload:
	./api-backend/scripts/initial-code-upload.sh

# Get a shell on the API backend Postgres pod (for debugging)
postgres-shell:
	bash -c "kubectl -n api exec -it masterdb-1-0 -- psql 'postgresql://postgres:$$(kubectl -n api get secret masterdb-secret -o jsonpath='{.data.superUserPassword}' | base64 --decode ; echo)@localhost/podinate'"

# After API spec change, rebuild the generate code
api-generate:
	bash api/generate.sh

salt-apply:
	ssh ubuntu@salt.podinate.com "rm -rf ~/salt/*"
	scp -r infrastructure/control/salt/* ubuntu@salt.podinate.com:~/salt/
	ssh ubuntu@salt.podinate.com "sudo cp -r ~/salt/* /srv/salt/ ; sudo salt '*' state.apply"