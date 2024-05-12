# Build stage
FROM golang:1.21.0 AS build-stage

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server cmd/api/main.go


# Deploy the application binary
FROM alpine:3.13
WORKDIR /app

ENV DB_NAME=
ENV DB_PORT=
ENV DB_HOST=
ENV DB_USERNAME=
ENV DB_PASSWORD=
ENV DB_PARAMS=
ENV JWT_SECRET=
ENV BCRYPT_SALT=

COPY --from=build-stage /app/bin/server .

EXPOSE 8080
CMD ["./server"]