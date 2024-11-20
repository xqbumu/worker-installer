.PHONY: dev
dev:
	wrangler dev

.PHONY: generate
generate:
	go generate ./...

.PHONY: build-wasm
build-wasm: generate
	go build -ldflags "-s -w" -o ./build/app.wasm .

.PHONY: build-tinygo
build-tinygo: generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	tinygo build -target wasm -o ./build/app.wasm -no-debug -panic=trap -gc=leaking

.PHONY: build-tinygo-wasip1
build-tinygo-wasip1: generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	GOOS=wasip1 GOARCH=wasm tinygo build -o ./build/app.wasm -no-debug -panic=trap -gc=leaking

.PHONY: build-go-js
build-go-js: generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm

.PHONY: build-go-wasip1
build-go-wasip1: generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	GOOS=wasip1 GOARCH=wasm go build -o ./build/app.wasm

.PHONY: install-qtc
install-qtc:
	go install github.com/valyala/quicktemplate/qtc@latest

.PHONY: deploy
deploy:
	wrangler deploy