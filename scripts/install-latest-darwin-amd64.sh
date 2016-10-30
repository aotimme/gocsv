#!/usr/bin/env bash

BINARY_DIRNAME="gocsv-darwin-amd64"
ZIP_FILENAME="${BINARY_DIRNAME}.zip"
ZIP_URL="https://github.com/DataFoxCo/gocsv/releases/download/latest/${ZIP_FILENAME}"
RAND=`date +%s | shasum | base64 | head -c 32; echo;`
mkdir /tmp/${RAND}
echo "Fetching binary..."
ZIP_ACTUAL_URL=`curl -s -I ${ZIP_URL} | grep "^Location: " | sed -n -e 's/^Location: //p' | tr -d '\r\n'`
curl -s ${ZIP_ACTUAL_URL} > /tmp/${RAND}/${ZIP_FILENAME}
echo "Extacting binary..."
unzip -q -d /tmp/${RAND} /tmp/${RAND}/${ZIP_FILENAME}
echo "Installing binary to /usr/local/bin/gocsv..."
mv /tmp/${RAND}/${BINARY_DIRNAME}/gocsv /usr/local/bin
echo "Cleaning up..."
rm -r /tmp/${RAND}
echo "GoCSV has been successfully installed!"
echo "Open a new Terminal window and run:"
echo "  gocsv help"
