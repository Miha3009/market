protoc --go_out=../inventory/internal/model --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative inventory.proto
