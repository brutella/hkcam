# Create a small(er) image from an sd-card.

**List disks**
```sh
# macOS
diskutil list
# Linux
fdisk -l
```

**Unmount disk on macOS (eg disk3)**
```sh
diskutil unmountDisk /dev/rdisk3
```

**Erase disk on macOS**
```sh
sudo diskutil eraseDisk FAT32 <name> MBRFormat /dev/disk3
# or
sudo diskutil zeroDisk /dev/disk3
```

**Install Raspbian on macOS**
```sh
sudo dd if=~/Downloads/2018-11-13-raspbian-stretch-lite.img of=/dev/rdisk3 bs=1m
```

**Create image from sd card on Linux**
```sh
sudo dd if=/dev/sda | gzip > image.img.gz bs=1M
```

**Copy disk image to sd-card on Linux**
```sh
gzip -dc image.img.gz | sudo dd of=/dev/sda bs=1M
```

# Create a custom Raspbian Stretch image

1. configure Raspberry Pi

- install [Raspbian](https://www.raspberrypi.org/downloads/raspbian/) 
- enable ssh
```sh
touch /Volumes/boot/ssh
```
- enable Wifi
```sh
# Define ssid and password
export WIFI_SSID=wifi; export WIFI_PWD=mypassword;

# Write network credentials into /Volumes/boot/wpa_supplicant.conf
echo "ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
ssid=\"$WIFI_SSID\"
psk=\"$WIFI_PWD\"
}" > /Volumes/boot/wpa_supplicant.conf
```

2. copy ssh key to the Raspberry Pi
```sh
ssh-copy-id pi@raspberrypi.local
```

3. run the `rpi` playbook
```sh
#! /bin/sh
cd $GOPATH/src/github.com/brutella/hkcam/ansible && ansible-playbook rpi.yml -i hosts
```
4. check if camera works
5. erase personal data from Raspberry Pi

- ssh on Raspberry Pi
```sh
#! /bin/sh
ssh pi@raspberrypi.local
```
- cleanup data
```sh
#! /bin/sh
sudo su
# stop services
sv stop hkcam

# delete hkcam data
rm -rf /var/lib/hkcam/data/*
rm -rf /var/log/hkcam/*

# delete wifi password
sh -c "echo 'ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
}' > /etc/wpa_supplicant/wpa_supplicant.conf"

# delete content from ansible and ssh-copy-id
rm -rf /home/pi/.ansible/
rm -rf ~/.ssh

# shutdown
shutdown now
```
6. put sd card into another linux machine
7. resize sd card
```sh
# shrink root file system to a minimum (-M)
e2fsck -f /dev/sda2
resize2fs -M /dev/sda2
#> The filesystem on /dev/sda2 is now 504923 (4k) blocks long.
# Remember the block size (4k) and count (504923)
# Now shrink the partition to 504923 * 4k = 2019692k
# https://askubuntu.com/questions/780284/shrinking-ext4-partition-on-command-line
fdisk /dev/sda
# 1. Delete partition 2
d
# 2. Create new primary partition 2
# - same START sector
# - new END sector +2019692K (note '+' and uppercase 'K')
n
# 3. Check partition table
p
# 3. Commit changes
w
# enlarge file system
resize2fs -p /dev/sda2
```
8. create disk image until last partition end – https://serverfault.com/a/853753
```sh
# Determine Units and End
fdisk -l -u=cylinders /dev/sda
# Example
#   Units: 2048 * 512 = 1048576 bytes
#   End:   2066
dd if=/dev/sda bs=1048576 count=2066 conv=sparse | gzip > image.img.gz
# or
dd if=/dev/sda of=~/image.img bs=1048576 count=2066 conv=sparse
``` 