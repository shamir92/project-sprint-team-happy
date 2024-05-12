# Build stage
FROM golang:1.21.0 AS build-stage

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server cmd/api/main.go


# Deploy the application binary
FROM alpine:3.13
WORKDIR /app

# COPY ./ap-southeast-1-bundle.pem .

COPY --from=build-stage /app/bin/server .

EXPOSE 8080
CMD ["./server"]