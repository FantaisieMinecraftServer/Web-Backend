FROM golang:latest

WORKDIR /
COPY . .

RUN go mod init backend \
  && go mod tidy \
  && go build

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

EXPOSE 8080


CMD ["go", "run", "main.go"]
