#!/bin/bash -e
FILENAME=pwndrop-linux-amd64
mkdir -p ${FILENAME}
cd ${FILENAME}
echo "*** downloading pwndrop."
wget https://github.com/kgretzky/pwndrop/releases/latest/download/${FILENAME}.tar.gz
echo "*** unpacking."
tar zxvf ${FILENAME}.tar.gz
cd pwndrop
chmod 700 pwndrop
echo "*** stopping pwndrop."
./pwndrop stop
echo "*** installing."
./pwndrop install
./pwndrop start
./pwndrop status
echo "*** cleaning up."
cd ../..
rm -rf ${FILENAME}/

