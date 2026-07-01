#!/usr/bin/env bash
set -e
SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.."; pwd)
go build "$ROOT_DIR/cmd/api"