#!/bin/bash
set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
rm -rf "$DIR/temp"
cp -Rf "$DIR/input" "$DIR/temp"
$DIR/../../vcgopkg "$DIR/temp" "-20060102150405"
unzip -l "$DIR/temp/go-detailedreport-to-csv/veracode/go-detailedreport-to-csv_cmd--go-detailedreport-to-csv--main-20060102150405.zip" | awk '{ print $4 }' | sort > $DIR/temp/temp.zip.log
unzip -l "$DIR/output/go-detailedreport-to-csv/veracode/go-detailedreport-to-csv_cmd--go-detailedreport-to-csv--main-20060102150405.zip" | awk '{ print $4 }' | sort > $DIR/temp/output.zip.log
diff $DIR/temp/output.zip.log $DIR/temp/temp.zip.log
rm -rf "$DIR/temp"