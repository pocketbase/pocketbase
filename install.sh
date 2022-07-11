#!/bin/sh

set -e

BINARY=pocketbase
INSTALL_DIR=/usr/local/bin
TAR_FORMAT=zip

ARCH=$(uname -m)
case $ARCH in
  x86) ARCH="386";;
  x86_64) ARCH="amd64";;
  i686) ARCH="386";;
  i386) ARCH="386";;
  arm64) ARCH="arm64";;
  *) echo "Unknown architecture: $ARCH"; exit 1;;
esac

OS=$(uname | tr '[:upper:]' '[:lower:]')

URL=$(curl -sSL https://api.github.com/repos/pocketbase/pocketbase/releases/latest | grep -o "https://.*${OS}_${ARCH}.${TAR_FORMAT}")

if [ -z "$URL" ]; then
    echo "No binary for your architecture and OS found. Please try again later."
    exit 1
fi

TARFILE=${BINARY}-${ARCH}-${OS}.${TAR_FORMAT}

temp_dir=$(mktemp -dt ${BINARY})
trap 'rm -rf $temp_dir' EXIT  INT TERM
cd "${temp_dir}"
curl --retry 10 -L -o "${TARFILE}" "${URL}"

echo "Installing to ${INSTALL_DIR}"

case $TAR_FORMAT in
  tar.gz) tar -xzf "${TARFILE}";;
  zip) unzip -qo "${TARFILE}";;
  *) echo "Unknown archive type: $TAR_FORMAT"; exit 1;;
esac

mv ${BINARY} ${INSTALL_DIR}
chmod +x ${INSTALL_DIR}/${BINARY}

echo "${BINARY} installed successfully!"
