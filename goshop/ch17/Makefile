# GoShop 的常用指令：make check、make build、make release
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)
PLATFORMS := linux/amd64 darwin/arm64 windows/amd64

.PHONY: check build release clean

# 提交前的檢查：格式、靜態分析、測試（含資料競爭偵測）
check:
	test -z "$$(gofmt -l .)"
	go vet ./...
	go test -race ./...

build:
	go build -ldflags "$(LDFLAGS)" -o bin/goshop .

# 跨平台編譯：GOOS／GOARCH 決定目標作業系統與 CPU 架構
release:
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=; \
		[ $$os = windows ] && ext=.exe; \
		echo "編譯 $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
			go build -ldflags "$(LDFLAGS)" -o bin/goshop-$$os-$$arch$$ext . ; \
	done

clean:
	rm -rf bin
