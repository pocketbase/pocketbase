# Build Golang binary
FROM golang:alpine AS builder-golang
WORKDIR /directory-to-build-golang-app
# Download dependencies in seperate step for docker layer caching
COPY [ "go.mod", "go.sum", "./" ]
RUN go mod download
# Copy everything else and build
COPY . . 
RUN cd examples/base && go build -o pocketbase ./

# Build the final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

# Copy pocketbase binary from the builder stage
COPY --from=builder-golang /directory-to-build-golang-app/examples/base/pocketbase /pb/pocketbase

# uncomment to copy the local pb_migrations dir into the image
# COPY ./pb_migrations /pb/pb_migrations

# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pb/pb_hooks

EXPOSE 8080

# start PocketBase
CMD ["/pb/pocketbase", "serve"]