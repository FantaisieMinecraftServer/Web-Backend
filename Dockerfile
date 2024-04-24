FROM golang:latest

WORKDIR /app
COPY ./app /app

RUN go mod init main \
  && go mod tidy \
  && go build

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

EXPOSE 6000


CMD ["go", "run", "app/main.go"]
