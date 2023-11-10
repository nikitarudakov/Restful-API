FROM golang:1.21.1

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /api ./cmd/app/main.go

EXPOSE 8080

CMD ["/api"]