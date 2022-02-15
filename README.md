# hkcam

`hkcam` is an open-source implementation of an HomeKit IP camera. 
It uses `ffmpeg` to access the camera stream and publishes the stream to HomeKit using [hc](https://github.com/brutella/hap).
The camera stream can be viewed in a HomeKit app. For example my [Home](https://hochgatterer.me/home) app works perfectly with `hkcam`.

## Features

- Live streaming via HomeKit
- Works with any HomeKit app
- [3D-Printed Enclosure](#enclosure)
- [Persistent Snapshots](#persistent-snapshots)
- Completely written in Go
- Runs on multiple platforms (Linux, macOS)

## Get Started

*hkcam uses Go modules and therefore requires Go 1.11 or higher.*

### Mac

The fastest way to get started is to

1. download the project on a Mac with a built-in iSight camera
```sh
git clone https://github.com/brutella/hkcam && cd hkcam
```
2. build and run `cmd/hkcam/main.go` by running `make run` in Terminal
3. open any HomeKit app and add the camera to HomeKit (pin for initial setup is `001 02 003`)

These steps require *git*, *go* and *ffmpeg* to be installed. On macOS you can install them via Homebrew.

```sh
brew install git
brew install go
brew install ffmpeg
```

### Raspberry Pi

If you want to create your own surveillance camera, you can run `hkcam` on a Raspberry Pi with attached camera module. 

#### Pre-configured Raspbian  Image

You can use a pre-configured Raspbian Stretch Lite image, where everything is already configured.

You only need to 

1. download the pre-configured Raspbian image and copy onto an sd card; [download](https://github.com/brutella/hkcam/releases/download/v0.0.9/raspbian-stretch-lite-2019-04-08-hkcam-v0.0.9-armv6.img.zip)

- **Note**: This image only works on a Raspberry Pi Zero

2. install [Etcher.app](https://www.balena.io/etcher/) and flash the downloaded image onto your sd card.
<img alt="Etcher.app" src="_img/etcher.png?raw=true"/>

> You can do the same on the command line as well.
> 
> On **macOS** you have to find the disk number for your sd card
> ```sh
> # find disk
> diskutil list
> ```
> You will see entries for `/dev/disk0`, `/dev/disk1`…, your sd card may have the disk number **3** and will be mounted at `/dev/disk3`
> 
> ```sh
> # unmount disk (eg disk3)
> diskutil unmountDisk /dev/rdisk3
> 
> # copy image on disk3
> sudo dd bs=1m if=~/Downloads/raspbian-stretch-lite-2019-04-08-hkcam-v0.0.9-armv6.img of=/dev/rdisk3 conv=sync
> ```

3. add your WiFi credentials so that the Raspberry Pi can connect to your WiFi

- create a new text file at `/Volumes/boot/wpa_supplicant.conf` with the following content
```sh
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
ssid="<ssid>"
psk="<password>"
}
```
- replace `<ssid>` with the name of your WiFi, and `<password>` with the WiFi password.
    
4. insert the sd card into your Raspberry Pi and power it up.
(After a reboot it may take up to several minutes until the camera is accessible via HomeKit – see [issue #136](https://github.com/brutella/hap/issues/136).)

5. open any HomeKit app and add the camera to HomeKit (pin for initial setup is `001 02 003`)


#### Manual Configuration

If you want, you can configure your Raspberry Pi manually.
This setup requires more configuration.
I've made an [Ansible](http://docs.ansible.com/ansible/index.html) playbook to configure your RPi with just one command.

The easiest way to get started is to

1. configure your Raspberry Pi

- install [Raspbian](https://www.raspberrypi.org/downloads/raspbian/) 
- [enable ssh](https://gist.github.com/brutella/0780479ceefc5d25a805b86ea795a3c6) (and WiFi if needed)
- connect a camera module

2. create ssh key and copy them to the Raspberry Pi
```sh
ssh-keygen
ssh-copy-id pi@raspberrypi.local
```

3 run the `rpi` playbook
```sh
cd ansible && ansible-playbook rpi.yml -i hosts
```

4. open any HomeKit app and add the camera to HomeKit (pin for initial setup is `001 02 003`)

These steps require *ansible* to be installed. On macOS you can install it via Homebrew.
```sh
brew install ansible
```

#### What does the playbook do?

The ansible playbook configures the Raspberry Pi in a way that is required by `hkcam`.
It does that by connecting to the RPi via ssh and running commands on it. 
You can do the same thing manually on the shell but ansible is more convenient.

Here are the things that the ansible playbook does.

1. Installs the required packages
    - [ffmpeg](http://ffmpeg.org) – to stream video from the camera via RTSP to HomeKit
    - [v4l2loopback](https://github.com/umlaeute/v4l2loopback) - to create a virtual video device to access the video stream by multiple ffmpeg processes
    - [runit](http://smarden.org/runit/) – to run `hkcam` as a service
2. Downloads and installs the latest `hkcam` release
3. Edits `/boot/config.txt` to enable access to the camera
4. Edits `/etc/modules` to enable the *bcm2835-v4l2* and *v4l2loopback* kernel modules
5. Restarts the RPi

After the playbook finishes, the RPi is ready to be used as a HomeKit camera.

**Additional Steps**

- I recommend to change the password of the `pi` user, once you have configured your Raspberry Pi.
- If you want to have multiple cameras on your network, you have to make sure that the hostnames are unqiue. By default the hostname of the Raspberry Pi is `raspberrypi.local`.
- SSH is enabled in the hkcam image. You may want to disable it.

**Debugging**

If experience issues with the hkcam daemon, you can find log outputs at `/var/log/hkcam/current`.

# Enclosure

<img alt="Desk mount" src="_img/enclosure-desk.jpg?raw=true" width="320" />
<img alt="Wall mount" src="_img/enclosure-wall.jpg?raw=true" width="320" />

The 3D-printed enclosure is designed for a Raspberry Pi Zero W and standard camera module.
You can use a stand to put the camera on a desk, or combine it with brackets of the [Articulating Raspberry Pi Camera Mount]() to mount it on a wall.

The 3D-printed parts are available as STL files [here](https://github.com/brutella/hkcam/tree/master/enclosure).

# Persistent Snapshots

In addition to video streaming, `hkcam` supports [Persistent Snapshots](/SNAPSHOTS.md).
*Persistent Snapshots* is a way to take snapshots of the camera and store them on disk.
You can then access them via HomeKit.

*Persistent Snapshots* are currently supported by [Home 3](https://hochgatterer.me/home), 
as you can see from the following screenshots.

| Services | Live Streaming | List of Snapshots |
|--------- | -------------- | ----------------- |
| <img alt="Services" src="_img/services.jpg?raw=true" width="280" /> | <img alt="Live streaming" src="_img/live-stream.jpg?raw=true" width="280" /> | <img alt="Snapshots" src="_img/snapshots.jpg?raw=true" width="280" /> |

| Snapshot | Automation | 
| --------------| -------------- |
| <img alt="Snapshot" src="_img/snapshot.jpg?raw=true" width="280" /> | <img alt="Automation" src="_img/automation.jpg?raw=true" width="280" /> |


# Advanced Configuration
The application can be further configured using flags in the startup script. These can lead to a misconfigured system and shoud be used at your own caution.

These settings can be changed in the startup script ```/etc/sv/hkcam/run```.

```
#!/bin/sh -e
exec 2>&1
v4l2-ctl --set-fmt-video=width=1280,height=720,pixelformat=YU12
exec hkcam --data_dir=/var/lib/hkcam/data --verbose=true
```

| Flag | Default value | Description |
|--------- | -------------- | ----------------- |
| min_video_bitrate | ```0``` | minimum video bit rate in kbps|
| multi_stream | ```false``` | "Allow mutliple clients to view the stream simultaneously|
| data_dir | ```"Camera"``` | Path to data directory|
| verbose | ```true```| Verbose logging|
| pin | ```"00102003"``` | PIN for HomeKit pairing |
| port | ```""``` | Port on which transport is reachable, random  portif empty |

## Network 
`hkcam` uses bonjour for service discovery. The port used for this ```5353```.
The transport port is random. It is assigned by the OS. You can set a port using the ```port``` flag.

# Contact

Matthias Hochgatterer

Website: [http://hochgatterer.me](http://hochgatterer.me)

Github: [https://github.com/brutella](https://github.com/brutella)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)


# License

`hkcam` is available under the Apache License 2.0 license. See the LICENSE file for more info.
