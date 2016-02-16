#!/bin/sh

# Copyright (C) 2014 Deepin Technology Co., Ltd.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 3 of the License, or
# (at your option) any later version.
# /etc/acpi/powerbtn.sh
# Initiates a shutdown when the power putton has been
# pressed.

[ -r /usr/share/acpi-support/power-funcs ] && . /usr/share/acpi-support/power-funcs

# getXuser gets the X user belonging to the display in $displaynum.
# If you want the foreground X user, use getXconsole!
getXuser() {
        user=`pinky -fw | awk '{ if ($2 == ":'$displaynum'" || $(NF) == ":'$displaynum'" ) { print $1; exit; } }'`
        if [ x"$user" = x"" ]; then
                startx=`pgrep -n startx`
                if [ x"$startx" != x"" ]; then
                        user=`ps -o user --no-headers $startx`
                fi
        fi
        if [ x"$user" != x"" ]; then
                userhome=`getent passwd $user | cut -d: -f6`
                export XAUTHORITY=$userhome/.Xauthority
        else
                export XAUTHORITY=""
        fi
        export XUSER=$user
}

# Skip if we just in the middle of resuming.
test -f /var/lock/acpisleep && exit 0

# If the current X console user is running a power management daemon that
# handles suspend/resume requests, let them handle policy This is effectively
# the same as 'acpi-support's '/usr/share/acpi-support/policy-funcs' file.

[ -r /usr/share/acpi-support/power-funcs ] && getXconsole
PMS="gnome-settings-daemon kpowersave xfce4-power-manager"
PMS="$PMS guidance-power-manager.py dalston-power-applet"
PMS="$PMS mate-settings-daemon power"

if pidof x $PMS > /dev/null; then
        exit
elif test "$XUSER" != "" && pidof dcopserver > /dev/null && test -x /usr/bin/dcop && /usr/bin/dcop --user $XUSER kded kded loadedModules | grep -q klaptopdaemon; then
        exit
elif test "$XUSER" != "" && test -x /usr/bin/qdbus; then
        kded4pid=$(pgrep -n -u $XUSER kded4)
        if test "$kded4pid" != ""; then
                dbusaddr=$(su - $XUSER -c "grep -z DBUS_SESSION_BUS_ADDRESS /proc/$kded4pid/environ")
                if test "$dbusaddr" != "" && su - $XUSER -c "export $dbusaddr; qdbus org.kde.kded" | grep -q powerdevil; then
                        exit
                fi
        fi
fi

# If all else failed, just initiate a plain shutdown.
/sbin/shutdown -h now "Power button pressed"
