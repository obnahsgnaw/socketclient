.PHONY: all
all: gateway

.PHONY: gateway
gateway:
	@ cd gateway/proto && buf generate
	@ echo "generated gateway proto"