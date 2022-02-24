grpc-generate:
	@cd proto && buf generate

grpc-lint:
	@cd proto && buf lint