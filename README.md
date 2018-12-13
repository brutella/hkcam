# hkcam

*hkcam* is an open-source implementation of an HomeKit IP-camera. It uses `ffmpeg` to access the camera stream and publishes the stream to HomeKit using [hc](https://github.com/brutella/hc). The camera stream can be viewed in a HomeKit app. For example my [Home](https//hochgatterer.me/home) app works perfectly with *hkcam*.

## Get Started

*hkcam uses Go modules and therefore requires Go 1.11+.*

### Mac

The fastest way to get started is to

1. download the project on a Mac with a built-in iSight camera

       # /bin/sh
       git clone https://github.com/brutella/hkcam && cd hkcam
    
2. build and run `cmd/hkcam/main.go` by running `make run` in Terminal
3. open any HomeKit app and add the camera to your home
4. enjoy your own HomeKit camera 🤓

> These steps require *git* and *go* to be installed. If you need to, you can install them via Homebrew.
> ```
> # /bin/sh
> brew install git
> brew install go
> ```

### RPi

If you want to create your own surveillance camera, you can run *hkcam* on a Raspberry Pi (~$25) with attached camera module (~20). This setup requires more configuration. I've made an [Ansible](http://docs.ansible.com/ansible/index.html) playbook to configure your RPi with just one command.

The easiest way to get started is to

1. install [Rasbian](https://www.raspberrypi.org/downloads/raspbian/) on your RPI and enable ssh (and WiFi if needed)
2. run the playbook
  
       # /bin/sh
       ansible-playbook rpi.yml -i raspberrypi.local,
    
3. open any HomeKit app and add the camera to your home
4. enjoy your own HomeKit camera 🤓

> This steps require *ansible* to be installed. If you're on a Mac you can install it via Homebrew.
> ```
> # /bin/sh
> brew install ansible
> ```

#### What does the playbook do?

The ansible playbook configures the Raspberry Pi in a way that is required by *hkcam*. It does that by connecting to the RPi via ssh and running commands on it. You can do the same thing manually on the shell but ansible does it automatically.

Here are the things that the ansible playbook does.

1. Installs the required packages
    - *ffmpeg* – to stream video from the camera via RTSP to HomeKit
    - [v4l2loopback](https://github.com/umlaeute/v4l2loopback) - to create a virtual video device to access the video stream by multiple ffmpeg processes
    - [runit](http://smarden.org/runit/) – to run *hkcam* as a service
2. Downloads and installs the latest *hkcam* release
3. Edits `/boot/config.txt` to enable access to the camera
4. Edits `/etc/modules` to enable the *bcm2835-v4l2* and *v4l2loopback* kernel modules
5. Restarts the RPi

After the playbook finishes, the RPi is ready to be used as a HomeKit camera.

**I highly recommend to change the password of the `pi` user once you haveconfigured your RPi.**

## Features

- 100% Go
- Works with any HomeKit app
- Runs on multiple platforms (Linux, macOS)

# Contact

Matthias Hochgatterer

Website: [http://hochgatterer.me](http://hochgatterer.me)

Github: [https://github.com/brutella](https://github.com/brutella)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)


# License

`hkcam` is available under the Apache License 2.0 license. See the LICENSE file for more info.
