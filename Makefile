PROTO := proto
PKG := pkg
# PROTO_TEMPLATES = $(shell find $(PKG) -iname "*.proto")
# PROTO_GEN = $(shell find $(PKG) -iname "*.pb.go")
PROTO_TEMPLATES = $(foreach dir, $(wildcard $(PKG)/*) , $(wildcard $(dir)/$(PROTO)/*.proto))
PROTO_GEN       = $(foreach dir, $(wildcard $(PKG)/*) , $(wildcard $(dir)/$(PROTO)/*.pb.go))

# .PHONY: build-protobuf
protobuf:
	@protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	$(PROTO_TEMPLATES)

clean-protobuf:
	@rm $(PROTO_GEN)
