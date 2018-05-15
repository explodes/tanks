DESKTOP_MAIN_GOPKG=github.com/explodes/tanks/go/cmd/app
DESKTOP_MAIN=./go/cmd/app/main.go
MOBILE_MAIN_GOPKG=github.com/explodes/tanks/go/cmd/mobile

ANDROID_PKG=io.explod.android
ANDROID_LIB=./android/tanklib/tankslib.aar
ANDROID_ACTIVITY=io.explod.android.minigames/io.explod.android.minigames.MainActivity

PROGRAM_NAME=tanks

BUILD_OUTPUT_DIR=./build

## +-+-+ STANDARD +-+-+ ##

clean-res:
	rm ./go/res/bindata.go || true

clean-build:
	find . -type d -name build | xargs -n1 rm -rf
	rm $(BUILD_OUTPUT_DIR) || true
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

create-generated-files: res

## +-+-+ GO DEPS +-+-+ ##

go-deps:
	go get -u -v github.com/jteeuwen/go-bindata/...
	go get -u -v golang.org/x/mobile/cmd/gomobile
	go get -u -v golang.org/x/mobile/cmd/gobind
	go get -u github.com/gopherjs/gopherjs

go-prep:
	gomobile init


## +-+-+ ANDROID +-+-+ ##

android-lib: create-generated-files
	CGO_ENABLED=1 gomobile bind -target android -javapkg "$(ANDROID_PKG)" -o "$(ANDROID_LIB)" "$(MOBILE_MAIN_GOPKG)"

# debug: build & run
android: enable-debug android-lib
	cd ./android; ./gradlew ':app:installDebug'
	mkdir -p "$(BUILD_OUTPUT_DIR)" || true
	cp ./android/app/build/outputs/apk/debug/app-debug.apk "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-android-debug.apk"
	adb shell am start -n "$(ANDROID_ACTIVITY)"

# debug: build
android-debug: enable-debug android-lib
	cd ./android; ./gradlew ':app:assembleDebug'
	mkdir -p "$(BUILD_OUTPUT_DIR)" || true
	cp ./android/app/build/outputs/apk/debug/app-debug.apk "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-android-debug-unsigned.apk"

# release: build
android-release: disable-debug android-lib
	cd ./android; ./gradlew ':app:assembleRelease'
	mkdir -p "$(BUILD_OUTPUT_DIR)" || true
	cp ./android/app/build/outputs/apk/release/app-release-unsigned.apk "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-android-release-unsigned.apk"

## +-+-+ WEB +-+-+ ##

# debug: build & run
web: enable-debug create-generated-files
	google-chrome http://localhost:8080
	CGO_ENABLED=1 gopherjs serve -m "$(DESKTOP_MAIN_GOPKG)"

# debug: build
web-debug: enable-debug create-generated-files
	CGO_ENABLED=1 gopherjs build -m -q -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-web-debug.js" "$(DESKTOP_MAIN_GOPKG)"

# release: build
web-release: disable-debug create-generated-files
	CGO_ENABLED=1 gopherjs build -m -q -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-web-release.js" "$(DESKTOP_MAIN_GOPKG)"

## +-+-+ LINUX +-+-+ ##

# pre-reqs:
# ubunutu: sudo apt install libglu1-mesa-dev libgles2-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libasound2-dev
# fedora: sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel alsa-lib-devel
# solus: sudo eopkg install libglu-devel libx11-devel libxrandr-devel libxinerama-devel libxcursor-devel libxi-devel

# debug: build & run
linux: enable-debug create-generated-files
	GOOS=linux GOARCH=amd64 go run "$(DESKTOP_MAIN)"

# debug: build
linux-debug: enable-debug create-generated-files
	GOOS=linux GOARCH=amd64 go build -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-linux-amd64-debug" "$(DESKTOP_MAIN_GOPKG)"

# release: build
linux-release: disable-debug create-generated-files
	GOOS=linux GOARCH=amd64 go build -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-linux-amd64-release" "$(DESKTOP_MAIN_GOPKG)"

## +-+-+ WINDOWS +-+-+ ##

# additional pre-reqs:
# ubunutu: sudo apt install gcc-multilib gcc-mingw-w64

# debug: build & run
windows: enable-debug create-generated-files
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc go run "$(DESKTOP_MAIN)"

# debug: build
windows-debug: enable-debug create-generated-files
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc go build -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-windows-i386-debug.exe" "$(DESKTOP_MAIN_GOPKG)"

# release: build
windows-release: disable-debug create-generated-files
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc go build -o "$(BUILD_OUTPUT_DIR)/$(PROGRAM_NAME)-windows-i386-release.exe" "$(DESKTOP_MAIN_GOPKG)"


## +-+-+ ALL +-+-+ ##

releases: android-release web-release linux-release windows-release
debugs: android-debug web-debug linux-debug windows-debug
all: releases debugs