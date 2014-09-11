#!/bin/bash

# Change log

# Version 1.0
# Date -----
# written by hippytaff
# Data aquisation phase written

# version 1.1
# Date 24/3/2011
# Written by hippytaff and matt_symes
# Diagnostic phase written

#>>> sudo rm /dev/rmkill ->  to remove an rfblock if 'sudo rfkill unblock all' does not work
#>>> rfkill unblock wifi
#>>> Add support for ndswrapper.
#>>> Add support for wicd
#>>> awk and grep to parse the files.
#>>> read each line from the file
#>>> parse line
#>>> structures in bash ? Arrays ?
#>>> Is rfkill installed  by default ?
#>>> wicd / wicd-daemon.py / ../wicd/daemon/monitor.py

################################################################################################
# Script description
################################################################################################

# Script arguments
# -verbose. Verbose output from continual monitoring.

################################################################################################
# script defines
################################################################################################

# Defines for the script.
VERBOSE=1
QUIET=0
WLESS=wireless-results.txt
VERSION=1.1
DATE_FN=
WHOAMI=
ROUTE=
PING=
NSLOOKUP=
WGET=
PGREP=
IFCONFIG=
DHCP=
PING_TRIES=3
REPING_TIME=0.2
ROUTER_PING_TARGET=	 	
EXTERNAL_PING_TARGET=8.8.8.8	# Ping multiple external targets in case one is down
DNS_LOOKUP_NAME=www.google.com	# Single target from multiple selections randomly.
WGET_URL=www.google.com		# Single target from multiple selections randomly. Use an array ?
NETWORK_MANAGER=[N]etworkManager
WPA_NAME=[w]pa_supplicant
PGREP_WICD_PID=
PGREP_NM_PID=
WICD=wicd
WLAN_NAME=
WLAN_DRIVER=
WLAN_IP=
RESOLV_FILE=/etc/resolv.conf
INTERFACES_FILE=/etc/network/interfaces
MODULES_FILE=/etc/modules
BLACKLIST_FILE=/etc/modprobe.d/blacklist.conf
BLACKLIST_FOLDER=/etc/modprobe.d/*
SYSLOG_LOG=/var/log/syslog
MESSAGES_LOG=/var/log/messages
KERNEL_LOG=/var/log/kern.log
NM_STATE_FILE=/var/lib/NetworkManager/NetworkManager.state
NM_APPLET_FILE=/etc/NetworkManager/nm-system-settings.conf

#################################################################################################
# script functions
#################################################################################################

# Function to check that rfkill is installed, and offer to install it if not
#
function call_rfkill_check
{
	RFKILL_INSTALLED=$(dpkg -l | grep -i rfkill)

	# Sleep for a second.
	sleep 1

	if [[ "$?" == "0" ]]
	then
		if  [[ "$RFKILL_INSTALLED" == "" ]]
		then
			# Ask the suer if they want to install rfkill
			echo "rfkill is not installed!...rfkill can unblock wireless blocks should there be any, would you like to install it?  y/n"
		
			# What do they say ?
			read ans
		
			if [[ "$ans" == "y" ]]
			then
				# might cause problems with non-debian distros. Consider removing the option to install
				#+or piss about to allow support for non-apt setups

				sudo apt-get -y install rfkill 	

				[[ "$?" -ne 0 ]] && { echo "Error installing rfkill"; return 1; }
			else
				echo "Continuing without rfkill."
				echo "You will not be notified of any blocks."

				return 1;
			fi
		fi
	else
		echo "Error checking for installed rfkill"

		return 1;
	fi

	sudo rfkill unblock all				# Now we know its there, unblock any blocks

	return 0;
}

# Function to redirect stdout to wireles-script.txt ($WLESS)
#
function call_redirect_stdout
{
	# Redirect output.
	exec 5>&1
	exec >> "$WLESS"
	exec 2>&1
}

# Function to restore stdout
#
function call_restore_stdout
{
	exec 1>&5 5>&-
}

# Function to check for that we havve root privilages
#
function call_check_root_privilages
{
	WHOAMI=$(which whoami)

	[[ "$WHOAMI" == "" ]] && { echo "Cannot find whoami binary"; exit 1; }

	# Check we are root.
	[[ $("$WHOAMI") == "root" ]] ||
		{ echo "This script needs to be run as root. Please run with sudo ./wireless_script"; exit 1; }
}

# Function to initialise the script
#
function call_initialise
{
	# Get the location of the 
	ROUTE=$(which route)
	PING=$(which ping)
	NSLOOKUP=$(which nslookup)
	WGET=$(which wget)
	PGREP=$(which pgrep)
	DATE_FN=$(which date)
	IFCONFIG=$(which ifconfig)
	DHCP=$(which dhclient)

	# Check binaries exist.
	# NOTE TO SELF. MAYBE SKIP WGET TEST IF BINARY NOT THERE.
	[[  "$ROUTE" != "" && "$PING" != "" &&
		"$NSLOOKUP" != "" && "$WGET" != "" && "$PGREP" != "" &&
		"$DATE_FN" != "" && "$IFCONFIG" != "" && "$DHCP" != "" ]] ||
		{ echo "cannot find required binaries"; exit 1; }

	sleep 1
}

# Function to check to see what networking daemons are running
#
function call_check_for_running_network_daemons
{
	echo "***********************************************************************************************"
	echo "Running networking services"
	echo "***********************************************************************************************"
	
	# Get the pid for network manager
	PGREP_NM_PID=$("$PGREP" -fx "$NETWORK_MANAGER")

	# Check for a running instance of network manager or networking. Could attempt to start if not running ?
	if [[ "$?" -eq 0 ]]
	then
		if [[ $PGREP_NM_PID != "" ]] 
		then
			echo "NetworkManager is running"

			return 0
		else
			echo "NetworkManager is _not_ running"
		fi
	else
		echo "NetworkManager is _not_ running"
	fi

	# No network manager ? Check for wicd.
	PGREP_WICD_PID=$("$PGREP" -fx "$WICD")

	if [[ "$?" -eq 0 ]]
	then
		if [[ $PGREP_WICD_PID != "" ]]
		then
			echo "WICD is running"

			return 0
		else
			echo "WICD is _not_ running"
		fi
	else
		echo "WICD is _not_ running"
	fi

	# Check to see that networking is up.

	return 0;
}

# Function to aquire the system info for the network
#
function call_get_system_info
{
	echo "************************************"
	echo "        Ubuntu release "
	echo "************************************"
	echo
	cat /etc/lsb-release
	echo
	echo "************************************"
	echo "        Kernel"
	echo "************************************"
	echo
	uname -a
	echo "************************************"
	echo "          List of drivers"
	echo "************************************"
	echo
	lsmod

	return 0
}

# Function to get the network info.
#
function call_get_network_info
{
	echo
	echo "************************************"
	echo "        pci wireless devices"
	echo "************************************"
	echo
	lspci -nnk | grep -i -A3 wirel
	echo
	echo "************************************"
	echo "        usb wireless devices"
	echo "************************************"
	echo
	lsusb | grep -i wirel
	echo
	echo "************************************"
	echo "        List of network devices"
	echo "************************************"
	echo
	sudo lshw -C Network
	echo
	sleep 3 #wait for 3 seconds
	echo
	echo "************************************"
	echo "           network info"
	echo "************************************"
	echo
	ifconfig -v -a
	echo
	echo "************************************"
	echo " Wireless specific network info"
	echo "************************************"
	echo
	iwconfig
	echo "************************************"
	echo " Rfkill Blocks"
	echo "************************************"
	echo
	sudo rfkill list all
	echo
	echo "************************************"

	return 0
}

# Function to parse
#
function call_parse_lshw_network
{
	# Get the network hardware details.
	LSHW_RES=$(sudo lshw -C Network)

	[[ "$?" -ne 0 ]] && { echo "Could not get lshw network information for parsing"; return 1; }

	# sleep for 3 seconds.
	sleep 3

	# Get the wireless interface name.
	WLAN_NAME=${LSHW_RES#*Wireless interface*logical name: }
#	WLAN_NAME=${WLAN_NAME%% *}
	WLAN_NAME=${WLAN_NAME%%$'\n'*}

	# Get the driver name
	WLAN_DRIVER=${LSHW_RES#*Wireless interface*driver=}
	WLAN_DRIVER=${WLAN_DRIVER%% *}

	# Get the IP address
	WLAN_IP=${LSHW_RES#*Wireless interface*ip=}
	WLAN_IP=${WLAN_IP%% *}

	# Sucess.
	return 0
}

# Function to get the file
#
function call_get_file_info
{
	function call_get_resolv
	{
		echo
		echo "**************************************************************************"
		echo "resolv.conf"
		echo "**************************************************************************"

		# Does the name server file exist
		[[ -f "$RESOLV_FILE" ]] && { cat "$RESOLV_FILE"; return 0; }

		echo "$RESOLV_FILE does not exist"

		return 1;
	}

	function call_get_interfaces
	{
		echo
		echo "*************************************************************************"
		echo "interfaces"
		echo "*************************************************************************"

		[[ -f "$INTERFACES_FILE" ]] && { cat "$INTERFACES_FILE"; return 0; }

		echo "$INTERFACES_FILE does not exist"

		return 1;
	}

	function call_get_blacklisted_devices
	{
		echo
		echo "*************************************************************************"
		echo "Blacklist file"
		echo "*************************************************************************"

		[[ -f "$BLACKLIST_FILE" ]] && { cat "$BLACKLIST_FILE"; return 0; }

		echo "$BLACKLIST_FILE does not exist"

		return 1;
	}

	function call_get_modules_file
	{
		echo
		echo "**************************************************************************"
		echo "Modules file"
		echo "**************************************************************************"

		[[ -f "$MODULES_FILE" ]] && { cat "$MODULES_FILE"; return 0; }

		echo "$MODULES_FILE does not exist"

		return 1;
	}

	function call_list_all_blacklist_files
	{
		echo
		echo "**************************************************************************"
		echo "Files in folder $BLACKLIST_FOLDER"
		echo "**************************************************************************"

		for blacklist_file in  $BLACKLIST_FOLDER
		do
			echo "$blacklist_file"
		done
	}

	function call_get_nm_state_file
	{
		echo
		echo "***************************************************************************"
		echo "NetworkManager.state"
		echo "***************************************************************************"

		[[ -f "$NM_STATE_FILE" ]] && { cat "$NM_STATE_FILE"; return 0; }

		echo "$NM_STATE_FILE does not exist"

		return 1;
	}

	function call_get_nm_applet_file
	{
		echo
		echo "****************************************************************************"
		echo "nm_applet_file"
		echo "****************************************************************************" 

		[[ -f "$NM_APPLET_FILE" ]] && { cat "$NM_APPLET_FILE"; return 0; }

		echo "$NM_APPLET_FILE does not exist"

		return 1;
	}

	# Get the files.
	call_get_interfaces
	call_get_resolv
	call_get_modules_file
	call_get_blacklisted_devices
	call_list_all_blacklist_files

	if [[ "$PGREP_NM_PID" != "" ]]
	then
		# Get nm specific files.
		call_get_nm_state_file
		call_get_nm_applet_file

	elif [[ "$PGREP_WICD_PID" != "" ]]
	then
		echo
	fi
}

# Function to parse the route
#
function call_parse_route
{
	# Get the routes. This will give us the ip address of the default gateway.
	# 0.0.0.0         xxx.xxx.xxx.xxx     0.0.0.0         UG    0      0        0 wlan0
	ROUTE_RES=$("$ROUTE" -n)

	# Did it fail ?
	[[ "$?" -eq 0  ]] || return 1

	echo
	echo "***************************************************************************"
	echo "Route info"
	echo "***************************************************************************"

	# Get the gateway ip addess. There nust be a better way to do this :(
	# A pass for each variable you need to get ?
	echo "$ROUTE_RES" | awk ' { print $0 } '
	ROUTER_PING_TARGET=$(echo "$ROUTE_RES" | awk ' $4 ~ /G/ { print $2 }')

	# Parse the route.
	[[ $ROUTER_PING_TARGET == "" ]] && { echo "Failed to parse route"; return 1; }

	# Sucess
	return 0
}

# Scan the access point.
#
function call_scan_AP
{
	echo
	echo "******************************************************************************"
	echo "Using nm-tool"
	echo "******************************************************************************"

	# Use nm-tool
	nm-tool

	echo
	echo "******************************************************************************"
	echo "Using iwlist scan"
	echo "******************************************************************************"

	# use iwlist
	iwlist scan
}

# Function to ping a target.
# $1 ping ip address.
# $2 ping count attempts
# $3 ping retry in seconds.
# $4 Verbose (1) or quiet
# $5 The interface to ping from
#
function call_ping_target
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" || "$3" == ""  || "$4" == ""  || "$5" == "" ]] && { echo "DEBUG: Interface value null in call_ping_target"; exit 1; }

	if (( $2 == $VERBOSE ))
	then
		echo
		echo "*********************************************************************************"
		echo "Ping test"
		echo "*********************************************************************************"
	fi

	# Ping and store the return string
	PING_RES=$("$PING" -c $2 -i $3 -I $5 $1)

	# What was the result of the ping.
	if [[ "$?" -eq 0 ]]
	then
		# success
		(( $4 == $VERBOSE )) && { echo "sucessfully pinged $1"; echo $PING_RES | awk ' { print $0 }'; }

		return 0
	else
		# failure
		echo "_unsucessfully_ pinged $1"

		return 1
	fi
}

# Function to perform a dns look up using nslookup
# $1 is the url to lookup
# $2 Verbose (1) or quiet (0)
#
function call_nslookup
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" ]] && { echo "DEBUG: Interface value null in call_nslookup"; exit 1; }

	if (( $2 == $VERBOSE ))
	then
		echo
		echo "*************************************************************************"
		echo "nslookup test"
		echo "*************************************************************************"
	fi

	# Perform a dns lookup on the required host.
	NS_LOOKUP_RES=$("$NSLOOKUP" $1)

	if [[ "$?" -eq 0 ]]
	then
		(( $2 == $VERBOSE )) && { echo "sucessfully looked up $1"; echo $NS_LOOKUP_RES | awk ' { print  $0 }'; }

		return 0
	else
		echo "_unsucessfully_ looked up $1"

		return 1
	fi
}

# Function to retrieve a file using wget
# $1 the url of the file to retrieve
# $2 Verbose (1)  or  Quiet (0)
#
function call_wget
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" ]] && { echo "DEBUG: Interface value null in call_wget"; exit 1; }

	if (( $2 == $VERBOSE ))
	then
		echo
		echo "*******************************************************************************"
		echo "wget test"
		echo "*******************************************************************************"
	fi

	WGET_RES=$("$WGET" -q "$1")

	if [[ "$?" -eq 0 ]]
	then
		(( $2 == $VERBOSE )) && { echo "sucessfully retrieved file $1"; echo $WGET_RES; }

		# Delete the file.
		rm "index.html"

		return 0
	else
		echo "_unsucessfully_ retrieved file $1"

		return 1
	fi
}

# Function to clean up, delete files etc
function call_cleanup
{
	echo
}

# Function to check thye connectivity to internal and external entities
# $1 Do we want continual probing ?
# $2 Interface name.
# $3 Override for verbosity. Passed into the script. CAN BE NULL.
#
function call_check_connectivity
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" ]] && { echo "DEBUG: Interface value null in call_check_connectivity"; exit 1; }

	echo
	echo "********************************************************************************"
	echo "Checking connectivity"
	echo "********************************************************************************"

	# Now we want to parse the network information
	call_parse_lshw_network

	# How did that pan out ?
	[[ "$?" -eq 0 ]] || { echo "Cannot perform connectivity checking"; return 1; }	

	MODE=$VERBOSE

	# loop	while :
	while :
	do
		# ping router. If ping fails interrogate dmesg. exit loop
		call_ping_target $ROUTER_PING_TARGET $PING_TRIES $REPING_TIME $MODE $2

		# Did it fail ?
		[[ "$?" == 0 ]] ||  { echo "Ping failure to router"; call_interrogate_logs; call_get_summary_network_info; return 1; }

		# ping external. If ping fails interrogate logs. exit loop
		call_ping_target $EXTERNAL_PING_TARGET $PING_TRIES $REPING_TIME $MODE $2

		# Did it fail ?
		[[ "$?" == 0 ]] || { echo "Ping failure to external server"; call_interrogate_logs; call_get_summary_network_info; return 1; }

		# Perform a dns lookup.
		call_nslookup $DNS_LOOKUP_NAME $MODE

		[[ "$?" == 0 ]] || { echo "nslookup failure"; call_interrogate_logs; call_get_summary_network_info; return 1; }

		# wget file. if wget fails interrogate logs. exit loop
		call_wget $WGET_URL $MODE

		# Did it fail ?
		[[ "$?" == 0 ]] || { echo "wget failure"; call_interrogate_logs; call_get_summary_network_info; return 1; }

		# Do we want continual probing ?
		[[ "$1" != "y" ]] && { return 0; }

		# sleep for a second. Hitting enter will  exit the loop.
		read -t 1 && return 0;

		# Change mode to quiet unless overridden by arguments. We don't want
		# the  log file to get too big.
		[[ $3 == "-verbose" ]] || { MODE=$QUIET; }
	done
}

function call_interrogate_logs
{
	# Interrogate all the logs we are interested in.
	call_interrogate_log "$SYSLOG_LOG"
	call_interrogate_log "$MESSAGES_LOG"
	call_interrogate_log "$KERNEL_LOG"
}

# Function to interrogate the logs for failure information.
# Currently interrogates /var/log/messages, /var/log/syslog and /var/log/kern.log
# $1 log to interrogate.
#
function call_interrogate_log
{
	# Sanity  check
	[[ "$1" == "" ]] && { echo "DEBUG: Interface value null in call_interrogate_log"; exit 1; }

	# Check the log file exists
	[[ -f "$1" ]] || { echo "Cannot find log file $1"; return 1; }

	# Log header.
	echo
	echo "******************************************************************************************"
	echo "******************************************************************************************"
	echo "Log file: $1"
	echo "******************************************************************************************"
	echo "******************************************************************************************"
	echo

	# Interrogate the log.
	cat "$1" | tail -n 30
}

# Function to restart network manager
# $1 Operation to be performed on Network Manager
#
function call_manage_nm
{
	# Sanity  check
	[[ "$1" == "" ]] && { echo "DEBUG: Interface value null in call_manager_nm"; exit 1; }

	echo "$1 Network Manager"
	echo

	# Restart network manager.
	sudo service network-manager "$1"

	# Sleep for 5 seconds.
	sleep 5
}

# Function to restart wicd
#
function call_restart_wicd
{
	echo "Restarting WICD"
	echo
}

# Function to remove and reload the kernel modules.
# $1 The module name
#
function call_reload_drivers
{
	# Sanity  check
	[[ "$1" == "" ]] && { echo "DEBUG: Interface value null in call_reload_drivers"; exit 1; }

	echo "Reloading driver $1"
	echo

	# Unload the driver
	$(sudo modprobe -r "$1")

	[[ "$?" -eq 0 ]] || { echo "Could not unload $1"; return 1; }

	# Reload the driver
	$(sudo modprobe "$1")

	[[ "$?" -eq 0 ]] || { echo "Could not reload $1"; return 1; }

	echo "Reloaded driver $1 sucessfully"

	# Sucess
	return 0;
}

# Function to drop and raise an interface
# $1 Interface name
#
function call_down_up_interface
{
	# Sanity  check
	[[ "$1" == "" ]] && { echo "DEBUG: Interface value null in call_down_up_interface"; exit 1; }

	echo "Dropping and raising interface $1"

	call_manipulate_interface "$1" "down"

	[[ "$?" -eq 0 ]] || { return 1; }

	sleep 3

	call_manipulate_interface "$1" "up"

	[[ "$?" -eq 0 ]] || { return 1; }

	echo "Dropped and raised interface $1 sucessfully"

	# Success
	return 0;
}

# Function to drop or raise an interface 
# $1 Interface
# $2 Operation to perform
#
function call_manipulate_interface
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" ]] && { echo "DEBUG: Interface value null in call_manipulate_interface"; exit 1; }

	$(sudo ifconfig "$1" "$2")

	[[ "$?" -eq 0 ]] || { echo "Could not "$2" interface $1"; return 1; }

	echo "$2 interface $1 sucessfully"

	return 0;
}

# Function to restart the networking service
#
function call_restart_networking
{
	sudo /etc/init.d/networking restart

	[[ "$?" -eq 0 ]] || { echo "Could not restart networking"; return 1; }

	return 0;
}

# A function to regressively try to connect to a router
# $1 Interface
# $2 Driver name.
#
function call_connect_regression
{
	# Sanity  check
	[[ "$1" == "" || "$2" == "" ]] && { echo "DEBUG: Interface value null in call_conect_via_script"; exit 1; }

	if [[ "$PGREP_NM_PID" != "" ]]
	then
		# Restart network manager
		call_manage_nm restart

	elif [[ "$PGREP_WICD_PID" != "" ]]
	then
		# Restart wicd.
		call_restart_wicd
	fi

	if [[ "$?" -eq 0 ]]
	then
		# Give time to try to reconnect
		sleep 60

		# Try to check the connectivity
		call_check_connectivity "n" "$1"

		[[ "$?" -eq 0 ]] && { echo "Restarting the manager restored connectivity "; return 0; }

	fi

	# Restart networking
	call_restart_networking
	
	if [[ "$?" -eq 0 ]]
	then

		# Give time to try to reconnect
		sleep 60

		# Try to check the connectivity
		call_check_connectivity "n" "$1"

		[[ "$?" -eq 0 ]] && { echo "Restarting the networking restored connectivity "; return 0; }
	fi

	# up down interface
	call_down_up_interface "$1"	

	if [[ "$?" -eq 0 ]]
	then
		# Give time to try to reconnect
		sleep 60

		# Try to check the connectivity
		call_check_connectivity "n" "$1"

		[[ "$?" -eq 0 ]] && { echo "Dropping and raising interface $1 restored connectivity "; return 0; }
	fi
	
	# Unload reload modules.
	call_reload_drivers "$2"

	if [[ "$?" -eq 0 ]]
	then
		# Give time to try to reconnect
		sleep 60

		# Try to check the connectivity
		call_check_connectivity "n" "$1"

		[[ "$?" -eq 0 ]] && { echo "Unloading and reloading $2 restored connectivity "; return 0; }
	fi

	echo "Failure. Could not restore connectivity."

	# Failure
	return 1;
}

# Function to get full system information from calls and files
#
function call_get_full_system_info
{
	# Get the system information
	call_get_system_info

	# Check for failure
	[[ "$?" -ne 0 ]] && return 1

	# Get the network info
	call_get_network_info

	# Check for failure
	[[ "$?" -ne 0 ]] && return 1

	# Get the file info.
	call_get_file_info

	# Parse the route to get the default gateway
	call_parse_route

	# Check for failure
	[[ "$?" -ne 0 ]] && return 1

	# Scan for access points.
	call_scan_AP

	return "$?";
}

# Function to get the basic summary frm system calls.
#
function call_get_summary_network_info
{
	# Get the basic system nformation
	call_get_network_info

	# Just pass on the previous return value.
	return "$?"
}

#####################################################################################################
# The meat and two veg ;)
#####################################################################################################

# clear the terminal screen.
clear

# Let the user know we are doing something
echo "Performing initial enviroment tests....."

# Delete any old file we many have
$(rm -rf wireless-results.txt)

# First things first. We neeed to be running with root privilages. 
call_check_root_privilages

# Initialise the script
call_initialise

# is rfkill installed
call_rfkill_check

# Set up the variable.
ans=

# Loop over getting the input.
while [[ "$ans" != [1-5] ]]
do
	# Clear screen again
	clear

	# Main menu
	#
	echo "This script will seek to return as much relevant information as possible to help diagnose the problem"
	echo
	echo "Please select from the following options"
	echo "1. Retrieve network status information"
	echo "2. Retrieve network status information and perform a single connectivity check"
	echo "3. Retrieve network status information and repeatedly perform connectivity checking"
	echo "	>>>>>> To exit continual connectivity checking hit the <enter> key <<<<<<"
	echo "4. Try to reconnect to disabled interface"
	echo "5. Exit"

	# Get user input.
	read ans
done

# If 6 was pressed exit the script here.
[[ "$ans" == "5" ]] && { echo "Good luck!"; exit 0; }

# Tell the user to be patisent.
# If they are smart they will put the kettle on at this point.
echo "Please be patient. Depending on the option you have selected this could take some time"

# redirect stdout
call_redirect_stdout

# Open code blocks for the forums.
echo "[code]"

# Initial log entries.
echo "Version: $VERSION (Development)"
echo $($DATE_FN)
echo

# Check to see what networking services are running.
call_check_for_running_network_daemons

case "$ans" in
	1)
		# At the the start get the full system information
		call_get_full_system_info
	;;
	2)
		# At the the start get the full system information
		call_get_full_system_info

		# Parse the network info
		call_parse_lshw_network

		# Did the last call fail for some reason ?
		[[ "$?" -eq 0 ]] && call_check_connectivity "n" "$WLAN_NAME" $1
	;;
	3)
		# At the the start get the full system information
		call_get_full_system_info

		# Parse the network info
		call_parse_lshw_network

		# Did the last call fail for some reason ?
		[[ "$?" -eq 0 ]] && call_check_connectivity "y" "$WLAN_NAME" $1
	;;
	4)
		# At the the start get the full system information
		call_get_full_system_info 

		# Parse the network info
		call_parse_lshw_network

		# Try to connect regresssively.
		[[ "$?" -eq 0 ]] && call_connect_regression "$WLAN_NAME" "$WLAN_DRIVER"
	;;

esac

# Finished
echo "Finished <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"

# Cleanup any mess we may have left.
call_cleanup

echo "http://wireless.kernel.org/en/users/Drivers"

# Close code blocks for the forums.
echo "[/code]"

#restore stdout
call_restore_stdout

echo "probe complete...please see 'wireless-results.txt'"
