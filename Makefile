.PHONY: dev
dev:
	wrangler dev

.PHONY: generate
generate:
	go generate ./...  

.PHONY: build
build: generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	# tinygo build -target wasm -o ./build/app.wasm
	tinygo build -target wasm -o ./build/app.wasm -no-debug -panic=trap
	# GOOS=js GOARCH=wasm go build -o ./build/app.wasm

.PHONY: deploy
deploy:
	wrangler deploy