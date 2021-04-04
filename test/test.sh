#!/bin/bash
for f in *; do
    if [ -d "$f" ]; then
        ./$f/test.sh
    fi
done