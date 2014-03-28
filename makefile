run:=go build -o .out && rm .out

all: 
	cd accounts/ &&  $(run)
	cd audio/ && $(run)
	cd datetime/ && $(run)
	cd desktop-toggle/ && $(run)
	cd display/ && $(run)
	cd dock-apps-builder && $(run)
	cd dock-daemon/ && $(run)
	cd grub2/ && $(run)
	cd inputdevices/ && $(run)
	cd keybinding/ && $(run)
	cd launcher-daemon/ && $(run)
	cd mime/ && $(run)
	cd mounts/ && $(run)
	cd network/ && $(run)
	cd power/ && $(run)
	cd system-info/ && $(run)
	cd themes/ && $(run)


update:
	sudo apt-get update && sudo apt-get install dde-go-dbus-factory go-dlib
