TANKS!
======

2-player tank game

Tanks race around screen attempting to lock on a good shot. Tanks fire automatically. Only the strongest will survive.

![Screenshot](https://raw.githubusercontent.com/explodes/go-wo/master/examples/tanks/tanks.png)

Running
-------

 - Android: `make android` (adb-ready connected Android device required)
 - Web: `make web` (opens up a browser tab)
 - Linux: `make linux` (starts the bloodshed pretty quickly, if you have the Linux Requirements)
 
Info
----

This is a port of [Tanks](https://github.com/explodes/go-wo/tree/master/examples/tanks) that written using 
[Pixel](https://github.com/faiface/pixel) but is now written with [Ebiten](https://github.com/hajimehoshi/ebiten) 
for portability.
 
Building release versions:
--------------------------

`make releases`

General requirements:
---------------------

 - [go-bindata](https://github.com/jteeuwen/go-bindata)
     - `go get -u github.com/jteeuwen/go-bindata/...`
 
Android requirements:
---------------------

 - [gomobile](https://github.com/golang/mobile)
     - `go get golang.org/x/mobile/cmd/gomobile`
     - `gomobile init # it might take a few minutes`
     
Web requirements:
-----------------

 - [gopherjs](https://github.com/gopherjs/gopherjs)
     - `go get -u github.com/gopherjs/gopherjs`
 

Linux requirements:
-------------------

ubunutu
```bash
sudo apt install libglu1-mesa-dev libgles2-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libasound2-dev
```

fedora
```bash
sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel alsa-lib-devel
```

solus
```bash
sudo eopkg install libglu-devel libx11-devel libxrandr-devel libxinerama-devel libxcursor-devel libxi-devel
```