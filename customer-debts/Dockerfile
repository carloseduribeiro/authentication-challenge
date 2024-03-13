FROM golang:1.22.0 AS builder

WORKDIR /app
COPY . .

RUN go get -v ./...
RUN GOOS=linux CGO_ENABLED=0 go build -C /app/cmd -ldflags="-w -s" -o runner main.go
RUN chmod +x /app/cmd/runner

FROM scratch

WORKDIR /app
COPY --from=builder /app/cmd/runner .
COPY ./cmd/.env .env

CMD ["/app/runner"]