#!/usr/bin/env bash

set -euo pipefail

curl -sfL "https://git.kernel.org/pub/scm/docs/man-pages/man-pages.git/plain/man5/proc.5" \
    | docker run -i --rm pandoc/core -f man -t plain --wrap=none \
    | python3 generate.py \
    | gofmt > generated.go
