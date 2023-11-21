FROM golang:1.21.1

WORKDIR /dao

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /api ./cmd/app/dao.go

EXPOSE 9000

CMD ["/api"]