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

# Set the desired timezone using the TZ environment variable
ENV TZ=Europe/Moscow

# Install the tzdata package to set the timezone in the container
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/yandex-food-service .

CMD ["./yandex-food-service"]
EXPOSE 80
