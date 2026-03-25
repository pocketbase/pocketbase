# Admin arayuzu (whitelabel) + ornek binary (examples/base)
# Derleme: docker build -t pocketbase-whitelabel .
# Urun adi: docker build --build-arg PB_PRODUCT_NAME="Benim Urun" -t pocketbase-whitelabel .

FROM node:22-alpine AS ui
WORKDIR /src/ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ ./
ARG PB_PRODUCT_NAME=PocketBase
ARG PB_LOGO_URL=
ENV PB_PRODUCT_NAME=${PB_PRODUCT_NAME}
ENV PB_LOGO_URL=${PB_LOGO_URL}
RUN npm run build

FROM golang:1.25-alpine AS app
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui /src/ui/dist ./ui/dist
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /pocketbase ./examples/base

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=app /pocketbase /usr/local/bin/pocketbase
WORKDIR /pb_data
EXPOSE 8090
VOLUME /pb_data
ENTRYPOINT ["/usr/local/bin/pocketbase"]
CMD ["serve", "--http=0.0.0.0:8090", "--dir=/pb_data"]
