FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-wallet ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/grpc-wallet ./
COPY --from=builder /app/.env ./

EXPOSE 50056

ENTRYPOINT ["./grpc-wallet"]