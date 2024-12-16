FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./migration ./migration

RUN go build -o /app/docker-wallet cmd/app/main.go


EXPOSE 8080

CMD ["/app/docker-wallet"]