.PHONY: all
all: proxy

.PHONY: proxy
proxy:
	@ cd go/proxy/proto && buf lint apis --error-format=json && buf generate --template buf.gen.yaml
	@ echo "generated proxy proto"