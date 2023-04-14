.PHONY: grpc
grpc:
	protoc --go_out=pkg/interfaces/grpc --go_opt=module=nononsensecode.com/grpc-tutorial/pkg/interfaces/grpc \
		--go-grpc_out=pkg/interfaces/grpc --go-grpc_opt=module=nononsensecode.com/grpc-tutorial/pkg/interfaces/grpc \
		api/grpc/*.proto