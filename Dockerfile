FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o ./cmd/run ./cmd/main.go

CMD ./cmd/run