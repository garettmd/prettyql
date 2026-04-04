#!/usr/bin/env bash
set -e

DIST_DIR="dist"

echo "Building WASM binary..."
mkdir -p "$DIST_DIR"
GOOS=js GOARCH=wasm go build -o "$DIST_DIR/prettyql.wasm" ./wasm/

echo "Copying wasm_exec.js from Go installation..."
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$DIST_DIR/wasm_exec.js"

echo "Copying frontend..."
cp web/index.html "$DIST_DIR/index.html"

# Get sizes for reporting
wasm_size=$(du -h "$DIST_DIR/prettyql.wasm" | cut -f1)
echo ""
echo "Build complete! Static site files in $DIST_DIR/:"
ls -lh "$DIST_DIR/"
echo ""
echo "To test locally:"
echo "  cd $DIST_DIR && python3 -m http.server 8080"
echo ""
echo "To deploy as a static site, upload the contents of $DIST_DIR/ to your host."
