# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
#COPY start.sh .
#COPY wait-for.sh .
#COPY database/migrations ./database/migrations

EXPOSE 7001
CMD [ "/app/main" ]
#ENTRYPOINT [ "/app/start.sh" ]
