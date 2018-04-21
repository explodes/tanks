

## +-+-+ STANDARD +-+-+ ##

clean-res:
	rm ./go/res/bindata.go || true

clean-build:
	find . -type d -name build | xargs -n1 rm -rf
	rm ./go/core/gen.go.bak || true
	rm ./android/tanklib/tankslib.aar || true
	rm ./android/tanklib/tankslib-sources.jar || true

clean: enable-debug clean-res clean-build

res: clean-res
	cd ./resources; $(GOPATH)/bin/go-bindata -nocompress -o ../go/res/bindata.go -pkg res ./...

enable-debug:
	sed -i.bak 's/Debug = false/Debug = true/g' ./go/core/gen.go

disable-debug:
	sed -i.bak 's/Debug = true/Debug = false/g' ./go/core/gen.go

## +-+-+ GO DEPS +-+-+ ##

go-deps:
	go get -u -v github.com/jteeuwen/go-bindata/...
	go get -u -v golang.org/x/mobile/cmd/gomobile
	go get -u -v golang.org/x/mobile/cmd/gobind
	go get -u github.com/gopherjs/gopherjs

go-prep:
	gomobile init


## +-+-+ ANDROID +-+-+ ##

android-lib: res
	CGO_ENABLED=1 gomobile bind -target android -javapkg io.explod.android -o ./android/tanklib/tankslib.aar github.com/explodes/tanks/go/cmd/mobile

# debug: build & run
android: enable-debug android-lib
	cd ./android; ./gradlew ':app:installDebug'
	mkdir -p ./build || true
	cp ./android/app/build/outputs/apk/debug/app-debug.apk ./build/tanks-debug.apk
	adb shell am start -n io.explod.android.minigames/io.explod.android.minigames.MainActivity

# release: build
android-release: disable-debug android-lib
	cd ./android; ./gradlew ':app:assembleRelease'
	mkdir -p ./build || true
	cp ./android/app/build/outputs/apk/release/app-release-unsigned.apk ./build/tanks-release-unsigned.apk

## +-+-+ WEB +-+-+ ##

# debug: build & run
web: enable-debug res
	google-chrome http://localhost:8080
	CGO_ENABLED=1 gopherjs serve -m github.com/explodes/tanks/go/cmd/app

# release: build
web-release: disable-debug res
	CGO_ENABLED=1 gopherjs build -m -q -o ./build/tanks.js github.com/explodes/tanks/go/cmd/app

## +-+-+ LINUX +-+-+ ##

# pre-reqs:
# ubunutu: sudo apt install libglu1-mesa-dev libgles2-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libasound2-dev
# fedora: sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel alsa-lib-devel
# solus: sudo eopkg install libglu-devel libx11-devel libxrandr-devel libxinerama-devel libxcursor-devel libxi-devel

# debug: build & run
linux: enable-debug res
	GOOS=linux GOARCH=amd64 go run ./go/cmd/app/main.go

# release: build
linux-release: disable-debug res
	GOOS=linux GOARCH=amd64 go build -o ./build/tanks github.com/explodes/tanks/go/cmd/app


## +-+-+ ALL +-+-+ ##

releases: android-release web-release linux-release