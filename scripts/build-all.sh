#!/usr/bin/env bash
set -euo pipefail

APP_NAME="${APP_NAME:-linkmngr}"
OUT_DIR="${OUT_DIR:-dist}"
VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

LDFLAGS="-s -w -X github.com/usama/linkmngr-cli/internal/cli.Version=${VERSION}"

targets=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "linux arm64"
  "windows amd64"
)

rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"

echo "Building ${APP_NAME} ${VERSION} (${COMMIT}) at ${DATE}"
echo

for target in "${targets[@]}"; do
  IFS=' ' read -r goos goarch <<< "${target}"

  ext=""
  if [[ "${goos}" == "windows" ]]; then
    ext=".exe"
  fi

  bin_name="${APP_NAME}_${VERSION}_${goos}_${goarch}${ext}"
  out_path="${OUT_DIR}/${bin_name}"

  echo "-> ${goos}/${goarch}"
  CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
    go build -trimpath -ldflags "${LDFLAGS}" -o "${out_path}" ./cmd/linkmngr

  (
    cd "${OUT_DIR}"
    sha256sum "${bin_name}" > "${bin_name}.sha256"
  )
done

echo
echo "Artifacts written to ${OUT_DIR}/"
