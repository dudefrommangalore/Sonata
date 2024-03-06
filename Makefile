# The Go programming language
GO=go

# The protobuf compiler
PROTOC=protoc

# The list of all protobuf files
PROTOS=$(shell find . -name '*.proto' -not -path './third_party/*')

# The list of generated Go files
GO_SRC=$(patsubst %.proto,%*.pb*.go,$(PROTOS))

# The protobuf plugin for Go
#GO_PLUGIN=--go_out=. --go-grpc_out=. --grpc-gateway_out=logtostderr=true:.
GO_PLUGIN=--go_out=. --go-grpc_out=require_unimplemented_servers=false:. --grpc-gateway_out=logtostderr=true:.

.PHONY: proto
proto: $(GO_SRC)

$(GO_OUT)%*.pb*.go: %.proto
	cd .. && $(PROTOC) -I . -I ./Presto.Sonata/third_party/proto $(GO_PLUGIN) Presto.Sonata/$<

.PHONY: clean
clean:
	$(GO) clean
	rm -rf $(GO_SRC)

.PHONY: search
search:
	go install ./cmd/search/... && \
		echo "search is installed at $GOPATH/bin/search"



