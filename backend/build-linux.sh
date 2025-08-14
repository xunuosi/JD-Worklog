#!/bin/bash
set -e

# --- 自动定位到脚本所在目录 ---
cd "$(dirname "$0")"

# --- 自动跳转到 go.mod 根目录 ---
GOMOD=$(go env GOMOD 2>/dev/null || true)
if [[ -z "$GOMOD" ]]; then
  echo "[ERR] go.mod not found. Run in module root. go mod init xxx"
  exit 1
fi
MODROOT=$(dirname "$GOMOD")
cd "$MODROOT"

# --- 配置 ---
APP_NAME="jd-worklog"
MAIN_PKG="./cmd/server"
CGO_ENABLED=0
CC_386="i686-linux-gnu-gcc"
CC_AMD64="x86_64-linux-gnu-gcc"
EXTRA_BUILD_TAGS=""

# --- 版本信息 ---
VERSION="${1:-dev}"
GIT_COMMIT=$(git rev-parse --short=8 HEAD 2>/dev/null || echo "nogit")
GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "none")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.Version=$VERSION -X main.GitCommit=$GIT_COMMIT -X main.GitTag=$GIT_TAG -X main.BuildTime=$BUILD_TIME"

# --- 预检 ---
if ! ls "$MAIN_PKG"/*.go >/dev/null 2>&1; then
  echo "[ERR] No .go files under $MAIN_PKG"
  exit 1
fi

# --- 输出目录 ---
DIST="dist"
mkdir -p "$DIST"

echo "----------------------------------------"
echo "Building $APP_NAME version $VERSION"
echo "Git: commit=$GIT_COMMIT tag=$GIT_TAG"
echo "Time: $BUILD_TIME"
echo "CGO_ENABLED=$CGO_ENABLED"
echo "Tags=$EXTRA_BUILD_TAGS"
echo "OutDir=$DIST"
echo "----------------------------------------"

build_one() {
  SUFFIX="$1"
  GOOS="$2"
  GOARCH="$3"
  CC="$4"
  OUT="$DIST/${APP_NAME}_${SUFFIX}"
  TAGS_OPT=""
  [[ -n "$EXTRA_BUILD_TAGS" ]] && TAGS_OPT="-tags \"$EXTRA_BUILD_TAGS\""

  echo "GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED CC=$CC"
  env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED CC=$CC \
    go build -trimpath -ldflags "$LDFLAGS" $TAGS_OPT -o "$OUT" "$MAIN_PKG"
  
  shasum -a 256 "$OUT" > "${OUT}.sha256.txt"
  echo "[OK] Build success: $OUT"
}

# ===== linux/386 =====
build_one "linux_386" "linux" "386" "$CC_386"

# ===== linux/amd64 =====
build_one "linux_amd64" "linux" "amd64" "$CC_AMD64"

# --- 打包 ZIP ---
ZIP_FILE="$DIST/${APP_NAME}_${VERSION}.zip"
zip -j "$ZIP_FILE" "$DIST"/* >/dev/null
echo "[Pack] Created ZIP archive: $ZIP_FILE"
