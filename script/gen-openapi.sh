#!/bin/bash

set -eux

CURRENT_SCRIPT=$(realpath "${BASH_SOURCE[0]}")
CURRENT_SCRIPT_DIR=$(dirname "$CURRENT_SCRIPT")
ROOT_DIR=$(dirname "$CURRENT_SCRIPT_DIR")

swag init --generalInfo ./cmd/api/main.go --output ./cmd/docs
