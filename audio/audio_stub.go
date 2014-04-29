package main

import "dlib/dbus"
import "fmt"

func (o *Audio) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio",
		"com.deepin.daemon.Audio",
	}
}

func (card *Card) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Card", card.index),
		"com.deepin.daemon.Audio.Card",
	}
}

func (sink *Sink) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Sink", sink.index),
		"com.deepin.daemon.Audio.Sink",
	}
}

func (source *Source) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Source", source.index),
		"com.deepin.daemon.Audio.Source",
	}
}

func (sinkInput *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Application", sinkInput.index),
		"com.deepin.daemon.Audio.Application",
	}
}

func (sourceOutput *SourceOutput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Application", sourceOutput.index),
		"com.deepin.daemon.Audio.Application",
	}
}
func (client *Client) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		fmt.Sprint("/com/deepin/daemon/Audio/Client", client.index),
		"com.deepin.daemon.Audio.Client",
	}
}

func (op *SinkInput) setPropName(key string, value interface{}) {
	switch key {
	case "Mute":
		if v, ok := value.(bool); ok && v != op.Mute {
			op.Mute = v
		}
	case "Volume":
		if v, ok := value.(uint32); ok && v != op.Volume {
			op.Volume = v
		}
	}
	dbus.NotifyChange(op, key)
}
