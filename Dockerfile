# Build stage
FROM golang:1.22.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate

COPY db/migration ./db/migration
COPY wait-for.sh .
COPY start.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
