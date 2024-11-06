.PHONY: all
all: gateway proxy

.PHONY: gateway
gateway:
	@ cd go/gateway/proto && buf generate
	@ echo "generated gateway proto"
.PHONY: proxy
proxy:
	@ cd go/proxy/proto && buf generate
	@ echo "generated proxy proto"