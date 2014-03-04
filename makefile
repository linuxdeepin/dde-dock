run:=go build -o .out && rm .out

all: 
	cd grub2/ && $(run)
	cd keybinding/ && $(run)
	cd accounts/ &&  $(run)
	cd datetime/ && $(run)
	cd mime/ && $(run)
	cd desktop/ && $(run)
	cd mounts/ && $(run)
	cd display/ && $(run)
	cd inputdevices/ && $(run)
	cd personalization/ && $(run)
	cd network/ && $(run)
	cd power/ && $(run)
	cd system-info/ && $(run)
	cd audio/ && $(run)
	cd launcher-daemon/ && $(run)


update:
	sudo apt-get update && sudo apt-get install dde-go-dbus-factory go-dlib
