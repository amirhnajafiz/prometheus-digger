#!/usr/bin/env bash
set -euo pipefail

sudo rm -f /usr/local/bin/promdigger
rm -rf ~/.promdigger

make clean
