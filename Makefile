dev-cluster:
	k3d cluster create podinate-dev --agents 2 -v "$PWD":/mnt/code/@agent:*
	alias k=kubectl
	k cluster-info