CFLAGS = -g -Wall -DLIBEXECDIR="$(/usr/local/libexec)"\
		 `pkg-config --cflags --libs glib-2.0 gtk+-3.0 x11 xext xtst xi gnome-desktop-3.0 upower-glib libnotify libcanberra-gtk3 gudev-1.0` -lm

gsd-power-manager:gsd-power-manager.c \
	gpm-common.c \
	gsd-backlight-helper.c \
	gsd-backlight-linux.c\
	gsd-power-plugin.c\
	gsd-input-helper.c\
	func.c\
	gsd-power-main.c

	$(CC) -o $@ $^ $(CFLAGS) 

.PHONY:clean

clean:
	rm *.o gsd-power-manager
