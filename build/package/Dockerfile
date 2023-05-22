# build environment
FROM golang:1.19-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

#ENV CGO_ENABLED=0
RUN go build -o /app/yandex-food-service ./cmd/yandex-food

# production environment
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/yandex-food-service .

CMD ["./yandex-food-service"]
EXPOSE 80
