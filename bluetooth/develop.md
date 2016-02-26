## Bluetooth module

### Bluez5 tookit

```
$ bluetoothctl

# list adapters
> list

# show default adapter information
> show

# set default adapter
> select 00:12:34:56:78:90

# power on adapter
> power on

# enable scan for devices
> scan on

# enbale agent
> agent on

# list devices
> deveices

# show device information
> info 00:00:12:34:56:78

# connect device
> connect 00:00:12:34:56:78

# pair device
> pair 00:00:12:34:56:78

# trust device
> trust 00:00:12:34:56:78

# remove device
> remove 00:00:12:34:56:78
```

### Debug Bluez5

Restart Bluez5 daemon in front-end
```
# systemctl stop bluetooth
# /usr/lib/bluez5/bluetooth/bluetoothd -n -d
```

For some cases, Bluez5 works weird, just remove all the profiles and
restart it
```
# rm -rf /var/lib/bluetooth/*
# systemctl restart bluetooth
```

### Debug PulseAudio

Restart pulseaudio and run in front-end
```
$ pulseaudio --kill; pulseaudio --start --daemonize=false -v
```

Use `pacmd` to list loaded modules
```
$ sudo apt-get install pulseaudio-utils
$ pacmd
> list-modules
```
