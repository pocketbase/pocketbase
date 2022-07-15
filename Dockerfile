# Build
# =====
FROM golang:1.18.4-alpine AS gobuilder

WORKDIR /app

COPY . .

# install git for version
RUN apk add --no-cache git

RUN cd examples/base \
  && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/pocketbase/pocketbase.Version=`git describe --tags --dirty --always`" -o ../../pocketbase


# Run
# ====
FROM scratch
WORKDIR /
COPY --from=gobuilder /app/pocketbase /pocketbase

VOLUME /pb_data
VOLUME /pb_public
EXPOSE 8090
CMD ["pocketbase", "serve", "--http", ":8090"]
