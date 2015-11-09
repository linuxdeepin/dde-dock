all:
	@echo need make target

gen-nm-const-enums:
	sed '/NM_SECRET_AGENT_CAPABILITY_LAST =/d' /usr/include/libnm/nm-dbus-interface.h > /tmp/nm-dbus-interface.h
	glib-mkenums --template nm_enum_gen.go.tpl /usr/include/libnm/nm-vpn-dbus-interface.h /usr/include/libnm/nm-setting.h /tmp/nm-dbus-interface.h > nm_enum_gen.go
	gofmt -w nm_enum_gen.go
