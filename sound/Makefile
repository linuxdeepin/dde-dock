CFLAGS=-g -Wall
LDFLAGS=-lpulse

libdde-pulse.so:dde-pulse.c dde-pulse.h
	$(CC) -fPIC -shared -o $@ $^ $(LDFLAGS) $(CFLAGS)
pulse:pulse.go dde-pulse.c
	$(MAKE) libdde-pulse.so
	go build pulse.go
clean:
	go clean
	rm *.so *.o
