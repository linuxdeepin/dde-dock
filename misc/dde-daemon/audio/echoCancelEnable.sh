#!/bin/bash
aecArgs="$*"
# If no "aec_args" are passed on to the script, use this "aec_args" as default:
[ -z "$aecArgs" ] && aecArgs="analog_gain_control=0 digital_gain_control=1"
newSourceName="echoCancelSource"

# Reload "module-echo-cancel"
echo Reload \"module-echo-cancel\" with \"aec_args=$aecArgs\"
pactl unload-module module-echo-cancel 2>/dev/null
if pactl load-module module-echo-cancel use_master_format=1 aec_method=webrtc aec_args=\"$aecArgs\" source_name=$newSourceName; then
	# Set a new default source and sink, if module-echo-cancel has loaded successfully.
	pacmd set-default-source $newSourceName
fi
