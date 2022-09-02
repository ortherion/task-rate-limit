GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_AUTH_NAME=auth_service
SERVICE_ANALYTICS_NAME=analytics_messaging
SERVICE_PATH=g6834/team17/api

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)
BUF_VERSION:="v0.56.0"

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.PHONY: run
run:
	go run cmd/auth/main.go

.PHONY: lint
lint:
	golangci-lint run ./...


# ----------------------------------------------------------------

.PHONY: build
build: deps generate .build

.build:
	go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-v -o ./bin/task-service$(shell go env GOEXE) ./cmd/task/main.go

# ----------------------------------------------------------------

.PHONY: generate
generate: .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	$(BUF_EXE) generate

.generate-finalize-go:
	mkdir -p pkg/$(SERVICE_AUTH_NAME) pkg/$(SERVICE_ANALYTICS_NAME)
	mv pkg/gitlab.com/$(SERVICE_PATH)/$(SERVICE_AUTH_NAME)/* pkg/$(SERVICE_AUTH_NAME)
	mv pkg/gitlab.com/$(SERVICE_PATH)/$(SERVICE_ANALYTICS_NAME)/* pkg/$(SERVICE_ANALYTICS_NAME)
	rm -rf pkg/gitlab.com/
	cd pkg/$(SERVICE_AUTH_NAME) && ls go.mod || (go mod init gitlab.com/$(SERVICE_PATH)/pkg/$(SERVICE_AUTH_NAME) && go mod tidy)
	cd pkg/$(SERVICE_ANALYTICS_NAME) && ls go.mod || (go mod init gitlab.com/$(SERVICE_PATH)/pkg/$(SERVICE_ANALYTICS_NAME) && go mod tidy)

# ----------------------------------------------------------------
.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

