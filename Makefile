LOCAL_BIN:=$(CURDIR)/bin

run:
	go run cmd/ocp-resource-api/main.go

lint:
	golint ./...

test:
	go test -v ./...

.PHONY: build
build: vendor-proto .generate .build

.PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-resource-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-resource-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-resource-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-resource-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-resource-api \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-resource-api/ocp-resource-api.proto
		mv pkg/ocp-resource-api/github.com/ozoncp/ocp-resource-api/api/ocp-resource-api/* pkg/ocp-resource-api/
		rm -rf pkg/ocp-resource-api/github.com
		mkdir -p cmd/ocp-resource-api

.PHONY: .build
.build:
		go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/ocp-resource-api ./cmd/ocp-resource-api/main.go

.PHONY: install
install: build .install

.PHONY: .install
install:
		go install cmd/grpc-server/main.go

.PHONY: vendor-proto
vendor-proto: .vendor-proto

.PHONY: .vendor-proto
.vendor-proto:
		@if [ ! -d vendor.protogen ]; then \
				mkdir -p vendor.protogen &&\
				mkdir -p vendor.protogen/api/ocp-resource-api;\
		fi

		cp api/ocp-resource-api/ocp-resource-api.proto vendor.protogen/api/ocp-resource-api;
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
		fi

.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u google.golang.org/grpc
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u github.com/envoyproxy/protoc-gen-validate
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install github.com/envoyproxy/protoc-gen-validate
		go get -u github.com/rs/zerolog