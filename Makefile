grpc-generate:
	@cd proto && buf generate

grpc-lint:
	@cd proto && buf lint

cluster:
	kind create cluster --name kainotomia || true

cluster-delete:
	kind delete cluster --name kainotomia

dev-deploy: cluster grpc-generate
	skaffold dev --port-forward=services

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest