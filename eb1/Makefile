
cleanres:
	rm ./eb/res/bindata.go || true

gores: cleanres
	cd ./eb/resources; $(GOPATH)/bin/go-bindata -o ../res/bindata.go -pkg res ./...

golib: gores
	gomobile bind -target android -javapkg io.explod.android -o ./mobile/android/game/game.aar github.com/explodes/scratch/ebgames/eb1/eb/mobile

android:
	cd ./mobile/android; ./gradlew ':app:buildRelease'

install: golib
	cd ./mobile/android; ./gradlew ':app:installDebug'
	adb shell am start -n io.explod.android.eb1/io.explod.android.eb1.MainActivity

install-only:
	cd ./mobile/android; ./gradlew ':app:installDebug'
	adb shell am start -n io.explod.android.eb1/io.explod.android.eb1.MainActivity