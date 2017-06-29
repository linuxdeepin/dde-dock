#!/bin/sh
case $1/$2 in
        pre/*)
        ;;
        post/*)
            gdbus call -y -d com.deepin.system.Power -o /com/deepin/system/Power -m com.deepin.system.Power.Refresh
        ;;
esac
