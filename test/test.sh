#!/bin/bash
set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
for f in $DIR/*; do
    if [ -d "$f" ]; then
        echo "$f/test.sh"
        $f/test.sh
    fi
done