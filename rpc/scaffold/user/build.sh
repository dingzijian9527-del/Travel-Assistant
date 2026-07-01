#!/usr/bin/env bash
RUN_NAME="user"
SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/../../.."; pwd)
OUTPUT_DIR="$SCRIPT_DIR/output"

mkdir -p "$OUTPUT_DIR/bin"
cp "$SCRIPT_DIR"/script/* "$OUTPUT_DIR"/
chmod +x "$OUTPUT_DIR/bootstrap.sh"

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o "$OUTPUT_DIR/bin/${RUN_NAME}" "$ROOT_DIR/cmd/user"
else
    cd "$ROOT_DIR"
    go test -c -covermode=set -o "$OUTPUT_DIR/bin/${RUN_NAME}" -coverpkg=./... ./cmd/user
fi
