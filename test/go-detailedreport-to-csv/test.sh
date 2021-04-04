#!/bin/bash
set -e
cp -Rf input temp
../../vcgopkg temp
diff -r input output
rm -rf temp