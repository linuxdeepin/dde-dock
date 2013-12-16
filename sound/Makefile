CFLAGS=-g -Wall
LDFLAGS=-lpulse

libdde-pulse.so:dde-pulse.c dde-pulse.h
	$(CC) -fPIC -shared -o $@ $^ $(LDFLAGS) $(CFLAGS)
audio:audio.go dde-pulse.c
	$(MAKE) libdde-pulse.so
	go build audio.go
clean:
	go clean
	rm *.so *.o
