# templa - Go templating utility
# See LICENSE file for copyright and license details.

.POSIX:

include config.mk

SRC = templa.go util.go

all: options templa

options:
	@echo templa build options:
	@echo LDFLAGS = $(LDFLAGS)

templa:
	go build -ldflags $(LDFLAGS)

clean:
	rm -f templa templa-$(VERSION).tar.gz

dist: clean
	mkdir -p templa-$(VERSION)
	cp -R LICENSE Makefile README.md config.mk $(SRC) go.mod go.sum templa.1 templa-$(VERSION)
	tar -cf templa-$(VERSION).tar templa-$(VERSION)
	gzip templa-$(VERSION).tar
	rm -rf templa-$(VERSION)

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f templa $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/templa
	mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	cp -f templa.1 $(DESTDIR)$(MANPREFIX)/man1
	chmod 644 $(DESTDIR)$(MANPREFIX)/man1/templa.1

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/templa $(DESTDIR)$(MANPREFIX)/man1/templa.1

.PHONY: all options templa clean dist install uninstall
