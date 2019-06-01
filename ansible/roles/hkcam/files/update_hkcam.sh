#!/bin/bash

CURRENT="$(/usr/bin/hkcam -version)"
LATEST="$(curl -s https://api.github.com/repos/brutella/hkcam/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')"

if [ $CURRENT != ${LATEST#?} ]; then
  echo "A newer version $LATEST is available."
  echo "Update in progress ..."
  wget -O hkcam-latest_linux_armhf.tar.gz https://github.com/brutella/hkcam/releases/download/${LATEST}/hkcam-${LATEST}_linux_armhf.tar.gz
  tar -xvf hkcam-latest_linux_armhf.tar.gz && rm -rf hkcam-latest_linux_armhf.tar.gz
  sudo sv stop hkcam
  sudo mv hkcam-linux_armhf/usr/bin/hkcam /usr/bin && rm -rf hkcam-latest_linux_armhf
  sudo sv start hkcam
  sudo sv status hkcam
  echo "Update to version $LATEST done."
else
  echo "No newer version is available. The latest version ($CURRENT) is already installed."
fi
