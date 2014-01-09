run:=go build -o .out && rm .out

all: 
	cd daccounts/ &&  $(run)
	cd date-time/ && $(run)
	cd default-app/ && $(run)
	cd desktop/ && $(run)
	cd display/ && $(run)
	cd ext-device/ && $(run)
	cd individuate/ && $(run)
	cd binding-manager/ && $(run)
	cd network/ && $(run)
	cd power/ && $(run)
	cd shutdown-manager/ && $(run)
	cd sound/ && $(run)
	cd system-info/ && $(run)


update:
	sudo apt-get update && sudo apt-get install dde-go-dbus-factory go-dlib
