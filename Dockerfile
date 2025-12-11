FROM golang:alpine3.23 AS builder

WORKDIR /app
COPY . .

WORKDIR /app/examples/base
RUN go build -o /build/pocketbase


FROM alpine:3.23 AS runtime

WORKDIR /app
COPY --from=builder /build/pocketbase /bin/pocketbase
ENTRYPOINT ["/bin/pocketbase"]