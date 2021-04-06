#!/bin/bash

set -eu

NAME="calui"
VERSION="0.2"
DEST="dist"
OS=(darwin darwin linux windows)
ARCH=(arm64 amd64 amd64 amd64)
LEN=$((${#OS[@]} - 1))

mkdir -p "$DEST"

for i in $(seq 0 "$LEN"); do
  GOOS=${OS[$i]}
  GOARCH=${ARCH[$i]}

  if [ "$GOOS" = "windows" ]; then
    OUT="$NAME.exe"
  else
    OUT="$NAME"
  fi
  
  echo "$GOOS $GOARCH -> ${OUT}"
  go build -o "$OUT"
  upx "$OUT" 
  tar zcf "$DEST/$NAME-$VERSION-${GOOS}_${GOARCH}.tar.gz" "$OUT"
  rm "$OUT"
done
