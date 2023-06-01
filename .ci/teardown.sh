#!/bin/bash

set -e -u -o pipefail

[ "${VERBOSE:-0}" -eq 0 ] || set -x

_log () { echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $1"; }

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd "$SCRIPT_DIR"

_log "removing existing instances..."
docker compose kill
docker compose rm -f