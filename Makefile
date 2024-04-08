all: build-local

build-local:
	go fmt ./...
	go mod download
	go build -o chat ./cmd/main.go

run: build-local
	./chat

build-docker:
	go fmt ./...
	docker build -t chat:latest .

docker-run:
	docker run -v ~/.aws:/root/.aws:ro -e AWS_PROFILE --rm -it chat:latest
