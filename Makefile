LDFLAGS="-X github.com/joyrex2001/kubedock/internal/config.Date=`date -u +%Y%m%d-%H%M%S`  \
		 -X github.com/joyrex2001/kubedock/internal/config.Build=`git rev-list -1 HEAD`   \
		 -X github.com/joyrex2001/kubedock/internal/config.Version=`git describe --tags`  \
		 -X github.com/joyrex2001/kubedock/internal/config.Image=joyrex2001/kubedock:`git describe --tags | cut -d- -f1`"

run:
	CGO_ENABLED=0 go run main.go server -P -v 2 --port-forward

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -tags containers_image_openpgp -o kubedock

gox:
	CGO_ENABLED=0 gox -os="linux darwin windows" -arch="amd64" \
		-output="dist/kubedock_`git describe --tags`_{{.OS}}_{{.Arch}}" -ldflags $(LDFLAGS)

docker:
	docker build . -t joyrex2001/kubedock:latest

clean:
	rm -f kubedock
	rm -rf dist
	go mod tidy
	rm -f coverage.out
	go clean -testcache

cloc:
	cloc --exclude-dir=vendor,node_modules,dist,_notes,_archive .

fmt:
	find ./internal -type f -name \*.go -exec gofmt -s -w {} \;
	go fmt ./...

test:
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 go test ./... -cover

lint:
	golint ./internal/...
	# errcheck ./internal/... ./cmd/...

cover:
	CGO_ENABLED=0 go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

deps:
	go install golang.org/x/lint/golint@latest
	go install github.com/kisielk/errcheck@latest
	go install github.com/mitchellh/gox@latest
	go install github.com/tcnksm/ghr@latest

.PHONY: run build gox docker clean cloc fmt test lint cover deps
