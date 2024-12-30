#!/bin/bash
set -e

platforms=(
  darwin-amd64
  darwin-arm64
  freebsd-386
  freebsd-amd64
  freebsd-arm64
  linux-386
  linux-amd64
  linux-arm
  linux-arm64
  windows-386
  windows-amd64
  windows-arm64
)

# get 1st argument as version
VERSION=$1

IFS=$'\n' read -d '' -r -a supported_platforms < <(go tool dist list) || true

for p in "${platforms[@]}"; do
  goos="${p%-*}"
  goarch="${p#*-}"
  if [[ " ${supported_platforms[*]} " != *" ${goos}/${goarch} "* ]]; then
    echo "warning: skipping unsupported platform $p" >&2
    continue
  fi
  ext=""
  if [ "$goos" = "windows" ]; then
    ext=".exe"
  fi
  cc=""
  cgo_enabled="${CGO_ENABLED:-0}"
  if [ "$goos" = "android" ]; then
    if [ "$goarch" = "amd64" ]; then
      cc="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android${ANDROID_SDK_VERSION}-clang"
      cgo_enabled="1"
    elif [ "$goarch" = "arm64" ]; then
      cc="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROID_SDK_VERSION}-clang"
      cgo_enabled="1"
    fi
  fi
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="$cgo_enabled" CC="$cc" go build -trimpath -ldflags="-s -w -X github.com/shufo/gh-pr-stats/cmd.Version=${VERSION}" -o "dist/${p}${ext}"
done
