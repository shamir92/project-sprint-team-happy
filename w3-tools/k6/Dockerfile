FROM --platform=linux/amd64 golang:1.22.3-alpine3.19 as builder

WORKDIR /app

COPY srv/go.mod srv/go.sum ./

RUN go mod download

COPY ./srv .

ENV CGO_ENABLED=1
RUN go build -o main .



FROM --platform=linux/amd64 ubuntu:24.04

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]