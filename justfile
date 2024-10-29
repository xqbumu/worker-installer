# Define variables
frontend_out := "build/app.wasm"

dev:
	wrangler dev
deploy:
	wrangler deploy
generate:
	go generate ./...  

build:
    just generate
    go build -ldflags "-s -w" -o {{frontend_out}} .


# Frontend recipe
t:
    just generate
    go run github.com/syumai/workers/cmd/workers-assets-gen@latest
    tinygo build -target wasm -o {{frontend_out}} -no-debug -panic=trap -gc=leaking


# Frontend recipe
tw:
    just generate
    go run github.com/syumai/workers/cmd/workers-assets-gen@latest
    GOOS=wasip1 GOARCH=wasm tinygo build -o {{frontend_out}} -no-debug -panic=trap -gc=leaking


# Frontend recipe
b:
    just generate
    go run github.com/syumai/workers/cmd/workers-assets-gen@latest
    GOOS=js GOARCH=wasm go build -o {{frontend_out}}


# Frontend recipe
w:
    just generate
    go run github.com/syumai/workers/cmd/workers-assets-gen@latest
    GOOS=wasip1 GOARCH=wasm go build -o {{frontend_out}}

# tinygo build -target wasm -o ./build/app1.wasm -no-debug -panic=trap  
# size: 1.1 MB / gzip 460 KB
# tinygo build -target wasm -o ./build/app2.wasm -no-debug -panic=trap -gc=leaking                    
# size: 959 KB / gzip 387 KB

i:
    go install github.com/valyala/quicktemplate/qtc@latest