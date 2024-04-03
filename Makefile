# Make a dev cluster
dev-cluster:
	echo "Creating k3d cluster podinate-dev"
	k3d cluster create podinate-dev
	kubectl config use-context k3d-podinate-dev
	echo "Creating database"
	kubectl create namespace podinate
	kubectl apply -f kubernetes/masterdb-postgres.yaml
	kubectl -n podinate rollout status --watch --timeout=180s statefulset/postgres
	sleep 10
	echo "Migrating database"
	make postgres-migrate
	kubectl apply -f kubernetes/controller.yaml -f kubernetes/controller-dev.yaml
	kubectl -n podinate rollout status --watch --timeout=180s deployment/podinate-controller
	./controller/scripts/initial-code-upload.sh
	kubectl -n podinate exec -it deployment/podinate-controller -- bash -c "cd /go/src/github.com/johncave/podinate/controller && go run ./ init --email someone@example.com --ip 127.0.0.1"
	mkdir -p testapp
	kubectl -n podinate cp $$(kubectl -n podinate get pod -l app=podinate-controller -o jsonpath='{.items[0].metadata.name}'):/profile.yaml testapp/credentials.yaml

# Install K3d on Arch Linux
install-dependencies:
	paru -S rancher-k3d-bin
	sudo pacman -S kubectl docker
	# curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
	# sudo install skaffold /usr/local/bin/
	# rm skaffold
	
# Show the logs for the API backend in the Kubernetes cluster while developing
dev-controller-logs:
	kubectl -n podinate logs -l app=podinate-controller -f

dev-controller-shell:
	kubectl -n podinate exec -it deployment/podinate-controller -- /bin/bash

dev-code-upload:
	./controller/scripts/initial-code-upload.sh
	kubycat ./kubycat.yaml

make dev-port-forward:
	kubectl -n podinate port-forward service/podinate-controller 31443:3000

# Get a shell on the API backend Postgres pod (for debugging)
postgres-shell:
	bash -c "kubectl -n podinate exec -it postgres-0 -- psql 'postgresql://postgres:$$(kubectl -n podinate get secret masterdb-secret -o jsonpath='{.data.superUserPassword}' | base64 --decode ; echo)@localhost/podinate'"

# Apply postgres migrations with atlas
postgres-migrate:
	kubectl apply -f kubernetes/atlas.yaml
	kubectl -n podinate rollout status --watch --timeout=180s statefulset/atlas
	kubectl -n podinate cp database/atlas.hcl atlas-0:/
	kubectl -n podinate exec -it atlas-0 -- /migrations/migration.sh
	kubectl -n podinate delete -f kubernetes/atlas.yaml
	
# After API spec change, rebuild the generate code
api-generate:
	bash api/generate.sh
