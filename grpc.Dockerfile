FROM golang:1.21.1

WORKDIR /grpc

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /api ./cmd/grpcServer/grpcServer.go

EXPOSE 9000

CMD ["/api"]