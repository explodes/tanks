TANKS!
======

2-player tank game

Tanks race around screen attempting to lock on a good shot. Tanks fire automatically. Only the strongest will survive.

![Screenshot](https://raw.githubusercontent.com/explodes/go-wo/master/examples/tanks/tanks.png)

Running
-------

 - Android: `make android`
 - Web: `make web`
 - Linux: `make linux`
 
Info
----

This is a port of [Tanks](https://github.com/explodes/go-wo/tree/master/examples/tanks) that written using 
[Pixel](https://github.com/faiface/pixel) but is now written with [Ebiten](https://github.com/hajimehoshi/ebiten) 
for portability.
 
Building release versions:
--------------------------

`make releases`

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