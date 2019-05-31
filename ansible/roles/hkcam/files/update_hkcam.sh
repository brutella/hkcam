#!/bin/bash

GET_URL="$(curl -s https://api.github.com/repos/brutella/hkcam/releases/latest | grep \"browser_download_url.*gz\" | cut -d '"' -f 4)"
LATEST="$(curl -s https://api.github.com/repos/brutella/hkcam/releases/latest | grep "browser_download_url.*gz" | cut -d '"' -f 4 | cut -d "/" -f 8)"
CURRENT="$(/usr/bin/hkcam -version)"
LATEST_CHECK=${LATEST#?}

if [ $LATEST_CHECK == $CURRENT ]; then
  echo "No need to update. The latest version is still installed"
  echo "Online latest is:" $LATEST_CHECK
  echo "Installed is:" $CURRENT
else
  echo "Update needed. This update will start now."
  echo $LATEST_CHECK
  echo $CURRENT
  wget $GET_URL
  tar -xvf hkcam-${LATEST}_linux_armhf.tar.gz
  rm -rf hkcam*.gz
  sudo sv stop hkcam
  sudo cp hkcam-${LATEST}_linux_armhf/usr/bin/hkcam /usr/bin
  rm -rf hkcam-*_linux_armhf
  sudo sv start hkcam
  rm update.sh
fi