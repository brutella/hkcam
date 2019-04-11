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

**Install Rasbian on macOS**
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

# Create a custom Rasbian Stretch image

1. configure Raspberry Pi

- install [Rasbian](https://www.raspberrypi.org/downloads/raspbian/) 
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

2. run the `rpi` playbook
```sh
#! /bin/sh
cd $GOPATH/src/github.com/brutella/hkcam/ansible && ansible-playbook rpi.yml -i hosts --ask-pass
#> SSH password: raspberry
```
3. check if camera works
4. erase personal data from Raspberry Pi

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
rm -rf /var/lib/hkcam
rm -rf /var/log/hkcam

# delete wifi password
sh -c "echo 'ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
}' > /etc/wpa_supplicant/wpa_supplicant.conf"

# delete content from rpi-source
rm -rf /root/linux*

# shutdown
shutdown now
```
5. put sd card into another linux machine
6. resize sd card
```sh
apt-get update && apt-get install parted
# shrink root file system to a minimum (-M)
e2fsck -f /dev/sda2
resize2fs -M /dev/sda2
#> The filesystem on /dev/sda2 is now 504923 (4k) blocks long.
# Remember the block size (4k) and count (504923)
# Now shrink the parition to 504923 * 4k = 2019692k
# https://askubuntu.com/questions/780284/shrinking-ext4-partition-on-command-line
fdisk /dev/sda
# 1. Delete partition 2
d
# 2. Create new primary partition 2
# - same START sector
# - new END sector +2019692K (note '+' and uppercase 'K')
n
# 3. Check parition table
p
# 3. Commit changes
w
# enlarge file system
resize2fs -p /dev/sda2
```
7. create disk image until last partition end – https://serverfault.com/a/853753
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