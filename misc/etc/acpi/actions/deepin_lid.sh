#! /bin/sh

# check if lid button exists
ls /proc/acpi/button/lid/* >/dev/null 2>&1 || exit 0

# enable external peripherals wakeup
handle_lid_open() {
	if grep XHC /proc/acpi/wakeup | grep -q disabled; then
		echo XHC > /proc/acpi/wakeup
	fi
}

# disable external peripherals wakeup
handle_lid_close() {
	if grep XHC /proc/acpi/wakeup | grep -q enabled; then
		echo XHC > /proc/acpi/wakeup
	fi
}

if grep -q open /proc/acpi/button/lid/*/state; then
	handle_lid_open
else
	handle_lid_close
fi
