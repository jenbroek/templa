# templa version
VERSION = 0.5.0

# paths
PREFIX = /usr/local
MANPREFIX = $(PREFIX)/share/man

# flags
LDFLAGS = '-w -s -X main.VERSION=$(VERSION)'
