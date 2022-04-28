grpc-generate:
	@cd proto && buf generate

grpc-lint:
	@cd proto && buf lint

cluster:
	k3d cluster create -c k3d.config.yml

cluster-delete:
	k3d cluster delete -c k3d.config.yml

dev-deploy:
	skaffold dev
