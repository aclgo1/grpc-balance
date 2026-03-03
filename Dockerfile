FROM golang:1.22-alpine AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-wallet ./cmd/main.go


FROM scratch
WORKDIR /app

COPY --from=builder /app/grpc-wallet ./

COPY --from=builder /app/.env ./ 

EXPOSE 50056

ENTRYPOINT ["./grpc-wallet"]