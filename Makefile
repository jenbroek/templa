.POSIX:

include config.mk

all: build

build:
	go build -ldflags $(LDFLAGS)

clean:
	go clean
	rm -f templa-$(VERSION).tar.gz

dist: clean
	install -Dt templa-$(VERSION) LICENSE README.md config.mk Makefile *.go go.mod go.sum templa.1
	tar czf templa-$(VERSION).tar.gz templa-$(VERSION)
	rm -rf templa-$(VERSION)

install: all
	install -Dt $(DESTDIR)$(PREFIX)/bin -m755 templa
	install -Dt $(DESTDIR)$(MANPREFIX)/man1 -m644 templa.1

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/templa $(DESTDIR)$(MANPREFIX)/man1/templa.1

.PHONY: all build clean dist install uninstall
