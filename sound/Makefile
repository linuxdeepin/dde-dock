CFLAGS=-g -Wall
LDFLAGS=-lpulse

audio:audio.go dde-pulse.c
	$(MAKE) libdde-pulse.so
	go build
libdde-pulse.so:dde-pulse.c dde-pulse.h
	$(CC) -fPIC -shared -o $@ $^ $(LDFLAGS) $(CFLAGS)
clean:
	go clean
	rm *.so *.o
