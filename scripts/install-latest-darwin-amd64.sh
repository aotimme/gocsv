#!/usr/bin/env bash

BINARY_DIRNAME="gocsv-darwin-amd64"
ZIP_FILENAME="${BINARY_DIRNAME}.zip"
ZIP_URL="https://github.com/DataFoxCo/gocsv/releases/download/latest/${ZIP_FILENAME}"
TMP_DIR=$(mktemp -d)
INSTALL_DIR="/usr/local/bin"
INSTALL_LOCATION="${INSTALL_DIR}/gocsv"

COLOR_GRAY=$(tput setaf 7)
COLOR_RED=$(tput setaf 1)
COLOR_GREEN=$(tput setaf 2)
TEXT_RESET=$(tput sgr0)

mark_fail() {
  echo -e "${COLOR_RED}\xE2\x9C\x98${TEXT_RESET}"
}

mark_pass() {
  echo -e "${COLOR_GREEN}\xE2\x9C\x94${TEXT_RESET}"
}

echo -n "Checking permissions... "
touch ${INSTALL_LOCATION}
if [ $? -eq 1 ]
then
  mark_fail
  echo "Unable to install gocsv."
  echo "You do not have permission to write to ${INSTALL_DIR}."
  echo "Change write permissions on that directory and try again."
  exit 1
fi
mark_pass

echo -n "Fetching binary... "
ZIP_ACTUAL_URL=`curl -s -I ${ZIP_URL} | grep "^Location: " | sed -n -e 's/^Location: //p' | tr -d '\r\n'`
curl -s ${ZIP_ACTUAL_URL} > ${TMP_DIR}/${ZIP_FILENAME}
mark_pass

echo -n "Extacting binary... "
unzip -q -d ${TMP_DIR} ${TMP_DIR}/${ZIP_FILENAME}
mark_pass

echo -n "Installing binary to ${INSTALL_LOCATION}... "
mv ${TMP_DIR}/${BINARY_DIRNAME}/gocsv ${INSTALL_LOCATION}
mark_pass

echo -n "Cleaning up... "
rm -r ${TMP_DIR}
mark_pass

echo ""
echo "GoCSV has been successfully installed!"
echo "Open a new Terminal window and run:"
echo "  gocsv help"
