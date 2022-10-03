# docker image just to build

FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app


# build proto

# Download proto zip
#ENV PROTOC_ZIP=protoc-3.15.8-linux-x86_64.zip
#RUN wget https://github.com/protocolbuffers/protobuf/releases/tag/v3.15.8/protoc-3.15.8-linux-x86_64.zip
#RUN unzip -o ${PROTOC_ZIP} -d ./proto 
#RUN chmod 755 -R ./proto/bin
#ENV BASE=/usr/local
# Copy into path
#RUN cp ./proto/bin/protoc ${BASE}/bin
#RUN cp -R ./proto/include/* ${BASE}/include

#RUN protoc --go_out=proto/generated --go_opt=paths=source_relative --go-grpc_out=proto/generated --go-grpc_opt=paths=source_relative proto/foodCard.proto
#protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative food.proto
# build app

ENV  GO111MODULE=on
ENV  CGO_ENABLED=0

RUN go build -o yandex-food-service ./cmd/api

RUN chmod +x /app/yandex-food-service

# small docker image to run 
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/yandex-food-service /app

CMD ["./app/yandex-food-service"]
EXPOSE 80
