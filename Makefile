.POSIX:

include config.mk

all: build

build:
	go build -ldflags $(LDFLAGS)

test:
	go test -coverpkg=./... ./...

fmt:
	find . -name '*.go' -exec gofmt -s -w '{}' +

clean:
	go clean
	rm -f templa-$(VERSION).tar.gz

dist: clean
	mkdir -p templa-$(VERSION)
	cp -R LICENSE README.md config.mk Makefile *.go internal go.mod go.sum templa.1 templa-$(VERSION)
	tar czf templa-$(VERSION).tar.gz templa-$(VERSION)
	rm -rf templa-$(VERSION)

install: all
	install -Dt $(DESTDIR)$(PREFIX)/bin -m755 templa
	install -Dt $(DESTDIR)$(MANPREFIX)/man1 -m644 templa.1

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/templa $(DESTDIR)$(MANPREFIX)/man1/templa.1

.PHONY: all build test fmt clean dist install uninstall
