FROM golang:latest

WORKDIR /cmd/backend
COPY ./cmd/backend /cmd/backend

RUN go mod init main \
    && go mod tidy \
    && go build

EXPOSE 8080


CMD ["go", "run", "backend.go"]
