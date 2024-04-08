FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Or -ldflags="-extldflags=-static"
ENV CGO_ENABLED=0
RUN go build -o chat ./cmd/main.go

FROM gcr.io/distroless/static

COPY --from=builder /app/chat /

ENTRYPOINT ["/chat"]

